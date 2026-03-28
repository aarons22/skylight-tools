package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-mcp/skylight/internal/client"
	"github.com/aarons22/skylight-mcp/skylight/internal/output"
	"github.com/spf13/cobra"
)

var (
	choresListAfter       string
	choresListBefore      string
	choresListIncludeLate bool
	choresListCategoryID  string
)

var choresListChoresCmd = &cobra.Command{
	Use:   "listChores",
	Short: "List chores with optional date range and category filters",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, _ := cmd.Root().PersistentFlags().GetString("base-url")
		configPath, _ := cmd.Root().PersistentFlags().GetString("config")
		frameIDFlag, _ := cmd.Root().PersistentFlags().GetString("frame-id")

		token, err := client.ResolveToken(configPath)
		if err != nil {
			return err
		}
		frameID, err := client.ResolveFrameID(frameIDFlag, configPath)
		if err != nil {
			return err
		}

		queryParams := map[string]string{}
		if choresListAfter != "" {
			queryParams["after"] = choresListAfter
		}
		if choresListBefore != "" {
			queryParams["before"] = choresListBefore
		}
		if choresListIncludeLate {
			queryParams["include_late"] = "true"
		}
		if choresListCategoryID != "" {
			queryParams["category_id"] = choresListCategoryID
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("GET", "/frames/{frameId}/chores",
			map[string]string{"frameId": frameID}, queryParams, nil)
		if err != nil {
			return err
		}

		jsonMode, _ := cmd.Root().PersistentFlags().GetBool("json")
		noColor, _ := cmd.Root().PersistentFlags().GetBool("no-color")
		if jsonMode {
			fmt.Fprintf(os.Stdout, "%s\n", string(resp))
			return nil
		}
		if err := output.PrintTable(resp, noColor); err != nil {
			fmt.Println(string(resp))
		}
		return nil
	},
}

func init() {
	choresCmd.AddCommand(choresListChoresCmd)
	choresListChoresCmd.Flags().StringVar(&choresListAfter, "after", "", "Start date filter (YYYY-MM-DD)")
	choresListChoresCmd.Flags().StringVar(&choresListBefore, "before", "", "End date filter (YYYY-MM-DD)")
	choresListChoresCmd.Flags().BoolVar(&choresListIncludeLate, "include-late", false, "Include overdue chores")
	choresListChoresCmd.Flags().StringVar(&choresListCategoryID, "category-id", "", "Filter by profile/category ID")
}
