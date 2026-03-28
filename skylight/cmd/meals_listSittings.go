package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/aarons22/skylight-tools/skylight/internal/output"
	"github.com/spf13/cobra"
)

var (
	mealsListSittingsDateMin string
	mealsListSittingsDateMax string
)

var mealsListSittingsCmd = &cobra.Command{
	Use:   "listSittings",
	Short: "List meal sittings in a date range",
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

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("GET", "/frames/{frameId}/meals/sittings",
			map[string]string{"frameId": frameID},
			map[string]string{
				"date_min": mealsListSittingsDateMin,
				"date_max": mealsListSittingsDateMax,
				"include":  "meal_category,meal_recipe",
			}, nil)
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
	mealsCmd.AddCommand(mealsListSittingsCmd)
	mealsListSittingsCmd.Flags().StringVar(&mealsListSittingsDateMin, "date-min", "", "Start date (YYYY-MM-DD)")
	mealsListSittingsCmd.Flags().StringVar(&mealsListSittingsDateMax, "date-max", "", "End date (YYYY-MM-DD)")
	_ = mealsListSittingsCmd.MarkFlagRequired("date-min")
	_ = mealsListSittingsCmd.MarkFlagRequired("date-max")
}
