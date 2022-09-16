package controllers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/cardboard-citizens/cz-goodboard-api/internal/discord"
	"github.com/cardboard-citizens/cz-goodboard-api/internal/models"
	"github.com/cardboard-citizens/cz-goodboard-api/internal/utils"
)

type Server struct {
	DB      *gorm.DB
	Discord *discordgo.Session
}

func (server *Server) InitializeDb(dbDriver, dbName string) (err error) {
	if dbDriver == "sqlite" {
		dbFile := fmt.Sprintf("%s.db", dbName)
		server.DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
		if err != nil {
			return err
		}
		utils.Log.Debug("Openned sqlite database connection on", dbFile)
	} else if dbDriver == "postgres" {
		server.DB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", dbName)), &gorm.Config{})
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Could not initialize database : Invalid or unsuported driver (%s)", dbDriver)
	}

	return err
}

func (server *Server) InitializeDiscord(botToken string) (err error) {
	server.Discord, err = discord.InitializeSession(botToken)
	return err
}

func (server *Server) MigrateDb() {
	server.DB.Debug().AutoMigrate(&models.Mission{})
}

func (server *Server) RegisterDiscordCommands(guildID string) {
	discord.RegisterCommands(server.Discord, guildID)
}
