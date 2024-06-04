package wallabag

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Entries struct {
	*WallabagClient
}

type PaginatedEntries struct {
	Page     int
	Limit    int
	Pages    int
	Total    int
	Embedded EmbeddedEntries `json:"_embedded"`
}

type EmbeddedEntries struct {
	Items []Entry `json:"items"`
}

type Entry struct {
	// base
	ID         int     `json:"id"`
	CreatedAt  *Time   `json:"created_at"`
	UpdatedAt  *Time   `json:"updated_at"`
	Title      *string `json:"title"`
	DomainName *string `json:"domain_name"`

	// content data
	Content        *string `json:"content"`
	Language       *string `json:"language"`
	ReadingTime    int     `json:"reading_time"`
	PreviewPicture *string `json:"preview_picture"`

	// article flags
	IsArchived IntBool `json:"is_archived"`
	ArchivedAt *Time   `json:"archived_at"`
	IsStarred  IntBool `json:"is_starred"`
	StarredAt  *Time   `json:"starred_at"`

	// tags
	Tags []Tag `json:"tags"`

	// annotations
	Annotations []Annotation `json:"annotations"`

	// urls
	URL            *string `json:"url"`
	HashedURL      *string `json:"hashed_url"`
	OriginURL      *string `json:"origin_url"`
	GivenURL       *string `json:"given_url"`
	HashedGivenURL *string `json:"hashed_given_url"`

	// published entries
	IsPublic    bool     `json:"is_public"`
	UID         *string  `json:"uid"`
	PublishedAt *Time    `json:"published_at"`
	PublishedBy []string `json:"published_by"`

	// user data
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`

	// fetching data
	Mimetype   *string           `json:"mimetype"`
	HttpStatus *string           `json:"http_status"`
	Headers    map[string]string `json:"headers"`
}

const (
	OptSortValueCreated  = "created"
	OptSortValueUpdated  = "updated"
	OptSortValueArchived = "archived"

	OptSortOrderAsc  = "asc"
	OptSortOrderDesc = "desc"

	OptDetailValueMetadata = "metadata"
	OptDetailValueFull     = "full"
)

func (e Entries) Exists(URLs ...string) (map[string]*int, error) {
	hashedURLs := make(map[string]string, len(URLs))
	hashes := make([]string, len(URLs))
	for i, it := range URLs {
		h := sha1.New()
		h.Write([]byte(it))
		v := hex.EncodeToString(h.Sum(nil))
		hashedURLs[v] = it
		hashes[i] = v
	}

	URL, _ := e.BuildURL("/api/entries/exists", nil)
	q := url.Values{"hashed_urls[]": hashes, "return_id": []string{"1"}}
	URL += "?" + q.Encode()

	resp, err := e.CallAPI(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	var raw map[string]*int
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	exists := make(map[string]*int, len(raw))
	for k, v := range raw {
		exists[hashedURLs[k]] = v
	}
	return exists, nil
}

type EntriesGetOptions struct {
	Archive    *bool
	Starred    *bool
	Sort       *string
	Order      *string
	Page       *int
	PerPage    *int
	Tags       []string
	Since      *time.Time
	Public     *bool
	Detail     *string
	DomainName *string
}

func (o *EntriesGetOptions) Validate() error {
	if o.Sort != nil {
		switch *o.Sort {
		case OptSortValueCreated, OptSortValueUpdated, OptSortValueArchived:
		default:
			return &InvalidOptionError{Field: "sort"}
		}
	}
	if o.Order != nil {
		switch *o.Order {
		case OptSortOrderAsc, OptSortOrderDesc:
		default:
			return &InvalidOptionError{Field: "order"}

		}
	}
	if o.Detail != nil {
		switch *o.Detail {
		case OptDetailValueMetadata, OptDetailValueFull:
		default:
			return &InvalidOptionError{Field: "detail"}
		}
	}
	return nil
}

func boolToStringedInt(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func (o *EntriesGetOptions) ToMap() map[string]string {
	data := make(map[string]string)
	if o.Archive != nil {
		data["archive"] = boolToStringedInt(*o.Archive)
	}
	if o.Starred != nil {
		data["starred"] = boolToStringedInt(*o.Starred)
	}
	if o.Sort != nil {
		data["sort"] = *o.Sort
	}
	if o.Order != nil {
		data["order"] = *o.Order
	}
	if o.Page != nil {
		data["page"] = strconv.Itoa(*o.Page)
	}
	if o.PerPage != nil {
		data["perPage"] = strconv.Itoa(*o.PerPage)
	}
	if o.Tags != nil {
		data["tags"] = strings.Join(o.Tags, ",")
	}
	if o.Since != nil {
		timestamp := o.Since.Unix()
		data["since"] = strconv.FormatInt(timestamp, 10)
	}
	if o.Public != nil {
		data["public"] = boolToStringedInt(*o.Public)
	}
	if o.Detail != nil {
		data["detail"] = *o.Detail
	}
	if o.DomainName != nil {
		data["domain_name"] = *o.DomainName
	}
	return data
}

func (e Entries) Get(options *EntriesGetOptions) (*PaginatedEntries, error) {
	url, err := e.BuildURL("/api/entries", options)
	if err != nil {
		return nil, err
	}
	resp, err := e.CallAPI(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	var entries PaginatedEntries
	start := time.Now()
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}
	log.Printf("Entries.Get parsed in %s", time.Since(start))
	return &entries, nil
}

type EntriesPostOptions struct {
	Title          *string
	Tags           []string
	Archive        *bool
	Starred        *bool
	Content        *string
	Language       *string
	PreviewPicture *string
	PublishedAt    *time.Time
	Authors        []string
	Public         *bool
	OriginURL      *string
}

func (o *EntriesPostOptions) Validate() error {
	return nil
}

func (o *EntriesPostOptions) ToMap() map[string]string {
	data := make(map[string]string)
	if o.Title != nil {
		data["title"] = *o.Title
	}
	if o.Tags != nil {
		data["tags"] = strings.Join(o.Tags, ",")
	}
	if o.Archive != nil {
		data["archive"] = boolToStringedInt(*o.Archive)
	}
	if o.Starred != nil {
		data["starred"] = boolToStringedInt(*o.Starred)
	}
	if o.Content != nil {
		data["content"] = *o.Content
	}
	if o.Language != nil {
		data["language"] = *o.Language
	}
	if o.PreviewPicture != nil {
		data["preview_picture"] = *o.PreviewPicture
	}
	if o.PublishedAt != nil {
		timestamp := o.PublishedAt.Unix()
		data["published_at"] = strconv.FormatInt(timestamp, 10)
	}
	if o.Authors != nil {
		data["authors"] = strings.Join(o.Authors, ",")
	}
	if o.Public != nil {
		data["public"] = boolToStringedInt(*o.Public)
	}
	if o.OriginURL != nil {
		data["origin_url"] = *o.OriginURL
	}
	return data
}

func (e Entries) Post(URL string, options *EntriesPostOptions) (*Entry, error) {
	var payload map[string]string
	if options == nil || reflect.ValueOf(options).IsNil() {
		payload = make(map[string]string)
	} else {
		payload = options.ToMap()
	}
	payload["url"] = URL

	url, err := e.BuildURL("/api/entries", nil)
	if err != nil {
		return nil, err
	}
	resp, err := e.CallAPI(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	var entry Entry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

func (e Entries) GetByID(entryID int) (*Entry, error) {
	URL, _ := e.BuildURL("/api/entries/"+strconv.Itoa(entryID), nil)
	resp, err := e.CallAPI(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	var entry Entry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// TODO test
func (e Entries) DeleteById(entryID int) (*Entry, error) {
	URL, _ := e.BuildURL("/api/entries/"+strconv.Itoa(entryID), nil)
	resp, err := e.CallAPI(http.MethodDelete, URL, nil)
	if err != nil {
		return nil, err
	}
	var entry Entry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}
	return &entry, nil
}
