package rest

import (
	"fmt"
	"github.com/cardboard-citizens/cz-mission-api/internal/database"
	"github.com/cardboard-citizens/cz-mission-api/internal/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

var (
	routeInitializers = make([]func(*RestController) error, 0)
)

type RestController struct {
	Port               string
	Router             *mux.Router
	DatabaseController *database.DatabaseController
}

func (controller *RestController) Initialize() (err error) {
	controller.Router = mux.NewRouter()

	for _, initializer := range routeInitializers {
		initializer(controller)
	}

	return err
}

func (controller *RestController) Listen() (err error) {
	corsConfig := cors.New(cors.Options{
		AllowedMethods:   []string{"POST", "PUT", "GET"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	})
	err = http.ListenAndServe(fmt.Sprintf(":%s", controller.Port), corsConfig.Handler(controller.Router))
	utils.Log.Info("REST server started on port", controller.Port)
	return err
}
