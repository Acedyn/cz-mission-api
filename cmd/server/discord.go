package main

import (
	"github.com/cardboard-citizens/cz-goodboard-api/internal/controllers"
	"github.com/spf13/cobra"
)

func deployDiscordCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "deploy-discord",
		Short: "Deploy discord commands",
		Long:  "Update the defined slash commands to the discord server",
		Run: func(cmd *cobra.Command, args []string) {
			server := controllers.Server{}
			server.InitializeDiscord(cmd.Flag("bottoken").Value.String())
			server.RegisterDiscordCommands(cmd.Flag("guildid").Value.String())
		},
	}

	command.Flags().StringP("bottoken", "b", "", "Private token to posess the bot")
	command.Flags().StringP("guildid", "g", "", "ID of the discord server")
	command.MarkFlagRequired("bottoken")
	command.MarkFlagRequired("guildid")
	return command
}

func init() {
	rootCmd.AddCommand(deployDiscordCommand())
}
