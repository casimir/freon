package cmd

import (
	"log"

	"github.com/casimir/freon/auth"
	"github.com/casimir/freon/database"
	"github.com/spf13/cobra"
)

var updatePasswordCmd = &cobra.Command{
	Use:   "update-password user password",
	Short: "Update the password of a user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		user, err := auth.FindUserByID(args[0])
		if err != nil {
			log.Fatalf("user not found: %s: %v", args[0], err)
		}
		if err := user.SetPassword(args[1]); err != nil {
			log.Fatalf("failed to generate password: %v", err)
		}
		if err := database.DB.Save(user).Error; err != nil {
			log.Fatalf("failed to save user: %v", err)
		}
	},
}
