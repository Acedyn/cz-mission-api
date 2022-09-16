package server

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/cardboard-citizens/cz-goodboard-api/internal/database"
	"github.com/cardboard-citizens/cz-goodboard-api/internal/discord"
	"github.com/cardboard-citizens/cz-goodboard-api/internal/utils"
)

type Server struct {
	initialized        bool
	DatabaseController *database.DatabaseController
	DiscordController  *discord.DiscordController
}

func (server *Server) Initialize(dbDriver, dbName, discordToken, discordGuildId string) (err error) {
	server.DatabaseController = &database.DatabaseController{
		DbDriver: dbDriver,
		DbName:   dbName,
	}
	err = server.DatabaseController.Initialize()
	if err != nil {
		return fmt.Errorf("Could not Initialize the server's db\n\t%s", err)
	}

	server.DiscordController = &discord.DiscordController{
		GuildId: discordGuildId,
	}
	err = server.DiscordController.Initialize(discordToken, server.DatabaseController)
	if err != nil {
		return fmt.Errorf("Could not Initialize discord connection\n\t%s", err)
	}

	err = server.DiscordController.RegisterCommands()
	if err != nil {
		return fmt.Errorf("Could not register discord commands\n\t%s", err)
	}

	server.initialized = true
	utils.Log.Info("Server initialization successfull")
	return err
}

func (server *Server) Run() (err error) {
	if !server.initialized {
		return fmt.Errorf("Could not run server: Initialization not completed")
	}

	server.DiscordController.ListenCommands()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	utils.Log.Info("Server running")
	fmt.Println("Press Ctrl+C to exit")
	<-stopSignal
	fmt.Println("Exit Captured for shutdown, press Ctrl+C again to force exit")

	err = server.Shutdown()
	if err != nil {
		return fmt.Errorf("Server did not shutdown gracefully\n\t%s", err)
	}
	return err
}

func (server *Server) Shutdown() (err error) {
	err = server.DiscordController.DeregisterCommands()
	if err != nil {
		return fmt.Errorf("Could not deregister discord commands\n\t%s", err)
	}

	err = server.DiscordController.Session.Close()
	if err != nil {
		return fmt.Errorf("Could not close discord session\n\t%s", err)
	}

	utils.Log.Info("Server shutdown successfull")
	return err
}
