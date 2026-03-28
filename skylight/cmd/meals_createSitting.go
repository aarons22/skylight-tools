package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/aarons22/skylight-tools/skylight/internal/output"
	"github.com/spf13/cobra"
)

var (
	mealsCreateSittingRecipeID         string
	mealsCreateSittingCategoryID       string
	mealsCreateSittingDate             string
	mealsCreateSittingAddToGroceryList bool
	mealsCreateSittingNote             string
	mealsCreateSittingRRule            string
	mealsCreateSittingDescription      string
)

var mealsCreateSittingCmd = &cobra.Command{
	Use:   "createSitting",
	Short: "Create a meal sitting",
	Long: `Create a meal sitting for a specific date.
Note: If --recipe-id is provided, do not set --description (API will return 422).`,
	Args: cobra.NoArgs,
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
			"meal_category_id":    mealsCreateSittingCategoryID,
			"date":                mealsCreateSittingDate,
			"add_to_grocery_list": mealsCreateSittingAddToGroceryList,
		}
		if mealsCreateSittingRecipeID != "" {
			body["meal_recipe_id"] = mealsCreateSittingRecipeID
		}
		if cmd.Flags().Changed("note") {
			body["note"] = mealsCreateSittingNote
		}
		if cmd.Flags().Changed("rrule") {
			body["rrule"] = mealsCreateSittingRRule
		}
		if cmd.Flags().Changed("description") {
			body["description"] = mealsCreateSittingDescription
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("POST", "/frames/{frameId}/meals/sittings",
			map[string]string{"frameId": frameID},
			map[string]string{
				"date_min": mealsCreateSittingDate,
				"date_max": mealsCreateSittingDate,
				"include":  "meal_category,meal_recipe",
			}, body)
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
	mealsCmd.AddCommand(mealsCreateSittingCmd)
	mealsCreateSittingCmd.Flags().StringVar(&mealsCreateSittingRecipeID, "recipe-id", "", "Meal recipe ID (optional)")
	mealsCreateSittingCmd.Flags().StringVar(&mealsCreateSittingCategoryID, "category-id", "", "Meal category ID (Breakfast/Lunch/Dinner/Snack)")
	mealsCreateSittingCmd.Flags().StringVar(&mealsCreateSittingDate, "date", "", "Sitting date (YYYY-MM-DD)")
	mealsCreateSittingCmd.Flags().BoolVar(&mealsCreateSittingAddToGroceryList, "add-to-grocery-list", false, "Add recipe ingredients to grocery list")
	mealsCreateSittingCmd.Flags().StringVar(&mealsCreateSittingNote, "note", "", "Optional note")
	mealsCreateSittingCmd.Flags().StringVar(&mealsCreateSittingRRule, "rrule", "", "RRULE for recurring sittings")
	mealsCreateSittingCmd.Flags().StringVar(&mealsCreateSittingDescription, "description", "", "Description (only when no recipe-id)")
	_ = mealsCreateSittingCmd.MarkFlagRequired("category-id")
	_ = mealsCreateSittingCmd.MarkFlagRequired("date")
}
