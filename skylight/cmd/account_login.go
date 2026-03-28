package cmd

import (
	"fmt"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/spf13/cobra"
)

var (
	accountLoginEmail    string
	accountLoginPassword string
	accountLoginFrameID  string
)

var accountLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and save credentials to config",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, _ := cmd.Root().PersistentFlags().GetString("base-url")
		configPath, _ := cmd.Root().PersistentFlags().GetString("config")

		token, err := client.Login(baseURL, accountLoginEmail, accountLoginPassword)
		if err != nil {
			return err
		}

		cfg := &client.Config{
			Token:   token,
			FrameID: accountLoginFrameID,
		}
		if err := client.SaveConfig(configPath, cfg); err != nil {
			return err
		}

		fmt.Println("Login successful. Token saved to config.")
		if accountLoginFrameID != "" {
			fmt.Printf("Frame ID saved: %s\n", accountLoginFrameID)
		}
		return nil
	},
}

func init() {
	accountCmd.AddCommand(accountLoginCmd)
	accountLoginCmd.Flags().StringVar(&accountLoginEmail, "email", "", "Skylight account email")
	accountLoginCmd.Flags().StringVar(&accountLoginPassword, "password", "", "Skylight account password")
	accountLoginCmd.Flags().StringVar(&accountLoginFrameID, "save-frame-id", "", "Optionally save a default frame ID to config")
	_ = accountLoginCmd.MarkFlagRequired("email")
	_ = accountLoginCmd.MarkFlagRequired("password")
}
