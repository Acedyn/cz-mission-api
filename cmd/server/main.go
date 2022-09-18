package main

import (
	"strconv"

	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cz-mission-api",
	Short: "API used to interact with CZ missions",
	Long:  `API that allow you to create and update missions from a discord API, and paricipate to these missions via REST. This API is dependend of the cz-auth-api package`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbosity, err := strconv.Atoi(cmd.Flag("verbosity").Value.String())
		if err != nil {
			utils.Log.Error(
				"Invalid verbosity level: Setting verbosity to DEBUG",
			)
		}
		utils.Log.SetLevel(verbosity)
	},
}

func init() {
	rootCmd.PersistentFlags().
		IntP("verbosity", "v", 20, "Verbosity level on a base of 10 (10 = DEBUG 50 = CRITICAL)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		utils.Log.Error(err)
	}
}
