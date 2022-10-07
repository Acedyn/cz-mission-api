package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/cardboard-citizens/cz-mission-api/internal/server"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
)

func serveCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Long:  "Start the database connection, REST API and the discord API",
		Run: func(cmd *cobra.Command, args []string) {
			server := server.Server{}
			err := server.Initialize(
				cmd.Flag("dbdriver").Value.String(),
				cmd.Flag("dbname").Value.String(),
				cmd.Flag("discordbottoken").Value.String(),
				cmd.Flag("discordguildid").Value.String(),
				cmd.Flag("port").Value.String(),
			)
			if err != nil {
				utils.Log.Error(err)
				return
			}
			err = server.Run()
			if err != nil {
				utils.Log.Error(err)
				return
			}
		},
	}

	command.Flags().StringP("port", "p", os.Getenv("HTTP_PORT"), "HTTP port to listen to")
	command.Flags().StringP("dbname", "n", os.Getenv("DB_NAME"), "Name of the database")
	command.Flags().
		StringP("dbdriver", "d", os.Getenv("DB_DRIVER"), "Driver that will be used to interact with the database (postgres, sqlite...)")
	command.Flags().
		StringP("discordbottoken", "b", os.Getenv("DISCORD_BOT_TOKEN"), "Private token to posess the bot")
	command.Flags().
		StringP("discordguildid", "g", os.Getenv("DISCORD_GUILD_ID"), "ID of the discord server")
	return command
}

func init() {
	godotenv.Load()
	rootCmd.AddCommand(serveCommand())
}
