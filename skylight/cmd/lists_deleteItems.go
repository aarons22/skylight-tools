package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aarons22/skylight-mcp/skylight/internal/client"
	"github.com/spf13/cobra"
)

var (
	listsDeleteItemsListID string
	listsDeleteItemsIDs    string
)

var listsDeleteItemsCmd = &cobra.Command{
	Use:   "deleteItems",
	Short: "Bulk delete list items (comma-separated IDs)",
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

		ids := strings.Split(listsDeleteItemsIDs, ",")
		for i := range ids {
			ids[i] = strings.TrimSpace(ids[i])
		}

		body := map[string]interface{}{
			"ids": ids,
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("DELETE", "/frames/{frameId}/lists/{listId}/list_items/bulk_destroy",
			map[string]string{"frameId": frameID, "listId": listsDeleteItemsListID}, nil, body)
		if err != nil {
			return err
		}

		jsonMode, _ := cmd.Root().PersistentFlags().GetBool("json")
		if jsonMode {
			fmt.Fprintf(os.Stdout, "%s\n", string(resp))
			return nil
		}
		fmt.Println("Items deleted successfully.")
		return nil
	},
}

func init() {
	listsCmd.AddCommand(listsDeleteItemsCmd)
	listsDeleteItemsCmd.Flags().StringVar(&listsDeleteItemsListID, "list-id", "", "List ID")
	listsDeleteItemsCmd.Flags().StringVar(&listsDeleteItemsIDs, "ids", "", "Comma-separated item IDs to delete")
	_ = listsDeleteItemsCmd.MarkFlagRequired("list-id")
	_ = listsDeleteItemsCmd.MarkFlagRequired("ids")
}
