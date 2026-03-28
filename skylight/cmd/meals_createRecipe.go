package cmd

import (
	"fmt"
	"os"

	"github.com/aarons22/skylight-tools/skylight/internal/client"
	"github.com/aarons22/skylight-tools/skylight/internal/output"
	"github.com/spf13/cobra"
)

var (
	mealsCreateRecipeCategoryID   string
	mealsCreateRecipeSummary      string
	mealsCreateRecipeDescription  string
)

var mealsCreateRecipeCmd = &cobra.Command{
	Use:   "createRecipe",
	Short: "Create a meal recipe",
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
			"meal_category_id": mealsCreateRecipeCategoryID,
			"summary":          mealsCreateRecipeSummary,
		}
		if cmd.Flags().Changed("description") {
			body["description"] = mealsCreateRecipeDescription
		}

		c := client.NewClient(baseURL, token)
		resp, err := c.Do("POST", "/frames/{frameId}/meals/recipes",
			map[string]string{"frameId": frameID},
			map[string]string{"include": "meal_category"}, body)
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
	mealsCmd.AddCommand(mealsCreateRecipeCmd)
	mealsCreateRecipeCmd.Flags().StringVar(&mealsCreateRecipeCategoryID, "category-id", "", "Meal category ID (Breakfast/Lunch/Dinner/Snack)")
	mealsCreateRecipeCmd.Flags().StringVar(&mealsCreateRecipeSummary, "summary", "", "Recipe name/title")
	mealsCreateRecipeCmd.Flags().StringVar(&mealsCreateRecipeDescription, "description", "", "Recipe description (ingredients and instructions)")
	_ = mealsCreateRecipeCmd.MarkFlagRequired("category-id")
	_ = mealsCreateRecipeCmd.MarkFlagRequired("summary")
}
