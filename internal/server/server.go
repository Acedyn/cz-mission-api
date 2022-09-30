package server

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/cardboard-citizens/cz-mission-api/internal/database"
	"github.com/cardboard-citizens/cz-mission-api/internal/discord"
	"github.com/cardboard-citizens/cz-mission-api/internal/rest"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
)

type Server struct {
	initialized        bool
	DatabaseController *database.DatabaseController
	DiscordController  *discord.DiscordController
	RestController     *rest.RestController
}

func (server *Server) Initialize(
	dbDriver, dbName, discordToken, discordGuildId string, port string,
) (err error) {
	server.DatabaseController = &database.DatabaseController{
		DbDriver: dbDriver,
		DbName:   dbName,
	}
	err = server.DatabaseController.Initialize()
	if err != nil {
		return fmt.Errorf("Could not Initialize the server's db\n\t%s", err)
	}

	server.DiscordController = &discord.DiscordController{
		GuildId:            discordGuildId,
		DatabaseController: server.DatabaseController,
	}
	err = server.DiscordController.Initialize(discordToken)
	if err != nil {
		return fmt.Errorf("Could not Initialize discord connection\n\t%s", err)
	}
	server.RestController = &rest.RestController{
		Port:               port,
		DatabaseController: server.DatabaseController,
	}
	err = server.RestController.Initialize()
	if err != nil {
		return fmt.Errorf("Could not Initialize rest router\n\t%s", err)
	}

	server.initialized = true
	utils.Log.Info("Server initialization successfull")
	return err
}

func (server *Server) Run() (err error) {
	if !server.initialized {
		return fmt.Errorf("Could not run server: Initialization not completed")
	}

	server.DiscordController.Listen()
	server.RestController.Listen()

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
	err = server.DiscordController.Shutdown()
	if err != nil {
		return fmt.Errorf("An error occured while shuting down the discord controller\n\t%s", err)
	}

	utils.Log.Info("Server shutdown successfull")
	return err
}
