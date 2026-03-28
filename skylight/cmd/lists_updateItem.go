package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/aarons22/skylight-tools/skylight/internal/output"
	"github.com/spf13/cobra"
)

var (
	listsUpdateItemListID  string
	listsUpdateItemItemID  string
	listsUpdateItemName    string
	listsUpdateItemChecked bool
	listsUpdateItemSetName bool
)

var listsUpdateItemCmd = &cobra.Command{
	Use:   "updateItem",
	Short: "Update a list item",
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

		attrs := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			attrs["name"] = listsUpdateItemName
		}
		if cmd.Flags().Changed("checked") {
			attrs["checked"] = listsUpdateItemChecked
		}

		body := map[string]interface{}{
			"data": map[string]interface{}{
				"type":       "list_items",
				"id":         listsUpdateItemItemID,
				"attributes": attrs,
			},
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("PATCH", "/frames/{frameId}/lists/{listId}/list_items/{itemId}",
			map[string]string{
				"frameId": frameID,
				"listId":  listsUpdateItemListID,
				"itemId":  listsUpdateItemItemID,
			}, nil, body)
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
	listsCmd.AddCommand(listsUpdateItemCmd)
	listsUpdateItemCmd.Flags().StringVar(&listsUpdateItemListID, "list-id", "", "List ID")
	listsUpdateItemCmd.Flags().StringVar(&listsUpdateItemItemID, "item-id", "", "Item ID")
	listsUpdateItemCmd.Flags().StringVar(&listsUpdateItemName, "name", "", "New item name")
	listsUpdateItemCmd.Flags().BoolVar(&listsUpdateItemChecked, "checked", false, "Checked state")
	_ = listsUpdateItemCmd.MarkFlagRequired("list-id")
	_ = listsUpdateItemCmd.MarkFlagRequired("item-id")
}
