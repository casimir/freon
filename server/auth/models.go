package auth

import (
	"errors"
	"log"
	"os"

	"github.com/casimir/freon/database"
	"github.com/casimir/freon/wallabag"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	database.DB.AutoMigrate(
		&User{},
		&WallabagCredentials{},
		&Token{},
		&TokenScope{},
	)

	// ensure at least one superuser exists
	result := database.DB.Select("id").Take(&User{}, "is_superuser = ?", true)
	notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if result.Error != nil && !notFound {
		panic(result.Error)
	}
	if notFound {
		log.Print("creating superuser...")
		username := os.Getenv("FREON_ADMIN_DEFAULT_USERNAME")
		if username == "" {
			username = "freon-admin"
		}
		user := User{
			Username:    username,
			IsSuperuser: true,
		}
		password := os.Getenv("FREON_ADMIN_DEFAULT_PASSWORD")
		if password == "" {
			password = "admin"
		}
		user.SetPassword(password)
		if result := database.DB.Create(&user); result.Error != nil {
			panic(result.Error)
		}
	}
}

type User struct {
	database.ModelUUID
	Username              string `gorm:"unique"`
	Password              []byte
	IsSuperuser           bool
	WallabagCredentialsID *uint
	WallabagCredentials   *WallabagCredentials
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err == nil {
		u.Password = hash
	}
	return err
}

func (u User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password)) == nil
}

func (u User) GetToken(ID string) (*Token, bool, error) {
	pk, err := uuid.Parse(ID)
	if err != nil {
		return nil, false, err
	}

	var token Token
	result := database.DB.Where("user_id = ?", u.ID).Take(&token, pk)
	if result.Error != nil {
		return nil, false, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, true, &UnknownModelError{"token", pk}
	}
	return &token, true, nil
}

func (u User) GetTokens() ([]Token, error) {
	var tokens []Token
	err := database.DB.Model(&Token{}).Where("user_id = ?", u.ID).Find(&tokens).Error
	return tokens, err
}

func (u *User) CreateToken(name string) error {
	token := Token{
		Name:   name,
		UserID: u.ID,
	}
	err := database.DB.Create(&token).Error
	return err
}

func (u *User) DeleteToken(ID string) (bool, error) {
	pk, err := uuid.Parse(ID)
	if err != nil {
		return false, err
	}

	result := database.DB.Delete(&Token{}, pk)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func FindUserByID(ID string) (*User, error) {
	pk, err := uuid.Parse(ID)
	if err != nil {
		return nil, err
	}

	var user User
	result := database.DB.Take(&user, pk)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, &UnknownUserError{ID}
	}
	return &user, nil
}

type WallabagCredentials struct {
	database.Model
	ServerURL     string          `json:"server_url"`
	ClientID      string          `json:"client_id" desc:"obfuscate"`
	ClientSecret  string          `json:"client_secret" desc:"obfuscate"`
	Username      string          `json:"username"`
	Password      string          `json:"password" desc:"obfuscate"`
	WallabagToken *wallabag.Token `gorm:"embedded" json:"-" desc:"hidden"`
}

func MustGetWallabagCredentials(ID uint) *WallabagCredentials {
	var creds WallabagCredentials
	result := database.DB.Take(&creds, ID)
	if result.Error != nil {
		panic(result.Error)
	}
	return &creds
}

func (w WallabagCredentials) ToCredentials() wallabag.Credentials {
	return wallabag.Credentials{
		ServerURL:    w.ServerURL,
		ClientID:     w.ClientID,
		ClientSecret: w.ClientSecret,
		Username:     w.Username,
		Password:     w.Password,
		Token:        w.WallabagToken,
	}
}

func (w *WallabagCredentials) UpdateWith(o *wallabag.Credentials) {
	w.ServerURL = o.ServerURL
	w.ClientID = o.ClientID
	w.ClientSecret = o.ClientSecret
	w.Username = o.Username
	w.Password = o.Password
	w.WallabagToken = o.Token
}

type Token struct {
	database.ModelUUID
	Name   string
	Scopes []TokenScope `gorm:"many2many:token_scopes;" desc:"hidden"`
	UserID uuid.UUID    `desc:"hidden"`
	User   User         `desc:"hidden"`
}

func (t *Token) UpdateWith(o *Token) {
	t.Name = o.Name
	t.Scopes = o.Scopes
}

type TokenScope struct {
	database.Model
	Name        string
	Description string
}
