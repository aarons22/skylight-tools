package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/aarons22/skylight-tools/skylight/internal/output"
	"github.com/spf13/cobra"
)

var (
	listsCreateItemListID  string
	listsCreateItemName    string
	listsCreateItemChecked bool
)

var listsCreateItemCmd = &cobra.Command{
	Use:   "createItem",
	Short: "Create a new list item",
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

		body := map[string]interface{}{
			"data": map[string]interface{}{
				"type": "list_items",
				"attributes": map[string]interface{}{
					"name":    listsCreateItemName,
					"checked": listsCreateItemChecked,
				},
			},
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("POST", "/frames/{frameId}/lists/{listId}/list_items",
			map[string]string{"frameId": frameID, "listId": listsCreateItemListID}, nil, body)
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
	listsCmd.AddCommand(listsCreateItemCmd)
	listsCreateItemCmd.Flags().StringVar(&listsCreateItemListID, "list-id", "", "List ID")
	listsCreateItemCmd.Flags().StringVar(&listsCreateItemName, "name", "", "Item name")
	listsCreateItemCmd.Flags().BoolVar(&listsCreateItemChecked, "checked", false, "Mark item as checked")
	_ = listsCreateItemCmd.MarkFlagRequired("list-id")
	_ = listsCreateItemCmd.MarkFlagRequired("name")
}
