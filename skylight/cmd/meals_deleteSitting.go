package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/aarons22/skylight-tools/skylight/internal/output"
	"github.com/spf13/cobra"
)

var (
	mealsDeleteSittingID      string
	mealsDeleteSittingDate    string
	mealsDeleteSittingDateMin string
	mealsDeleteSittingDateMax string
)

var mealsDeleteSittingCmd = &cobra.Command{
	Use:   "deleteSitting",
	Short: "Delete a meal sitting instance",
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

		queryParams := map[string]string{
			"include": "meal_category,meal_recipe",
		}
		if mealsDeleteSittingDateMin != "" {
			queryParams["date_min"] = mealsDeleteSittingDateMin
		}
		if mealsDeleteSittingDateMax != "" {
			queryParams["date_max"] = mealsDeleteSittingDateMax
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("DELETE",
			"/frames/{frameId}/meals/sittings/{sittingId}/instances/{date}",
			map[string]string{
				"frameId":   frameID,
				"sittingId": mealsDeleteSittingID,
				"date":      mealsDeleteSittingDate,
			}, queryParams, nil)
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
	mealsCmd.AddCommand(mealsDeleteSittingCmd)
	mealsDeleteSittingCmd.Flags().StringVar(&mealsDeleteSittingID, "sitting-id", "", "Sitting ID")
	mealsDeleteSittingCmd.Flags().StringVar(&mealsDeleteSittingDate, "date", "", "Instance date to delete (YYYY-MM-DD)")
	mealsDeleteSittingCmd.Flags().StringVar(&mealsDeleteSittingDateMin, "date-min", "", "Date range min for response (YYYY-MM-DD)")
	mealsDeleteSittingCmd.Flags().StringVar(&mealsDeleteSittingDateMax, "date-max", "", "Date range max for response (YYYY-MM-DD)")
	_ = mealsDeleteSittingCmd.MarkFlagRequired("sitting-id")
	_ = mealsDeleteSittingCmd.MarkFlagRequired("date")
}
