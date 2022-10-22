package rest

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/cardboard-citizens/cz-mission-api/internal/models"
)

func getMissionRoutes(controller *RestController) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/{id}/validate": func(w http.ResponseWriter, r *http.Request) {
			id := mux.Vars(r)["id"]
			user, err := controller.DatabaseController.GetOrCreateUserFromString(r.URL.Query().Get("user"))
			if err != nil {
				CustomResponse(w, "An error occured while initializing the user", []error{err}, http.StatusInternalServerError)
				return
			}
			mission, err := controller.DatabaseController.GetMissionFromString(id)
			if err != nil || mission == nil {
				CustomResponse(w, "Could not get mission", []error{}, http.StatusNotFound)
				return
			}

			participation, err := controller.DatabaseController.ValidateMission(mission, user)
			if err != nil {
				CustomResponse(w, "Could not validate mission", []error{}, http.StatusNotFound)
				return
			}

			SuccessResponse(w, participation, make([]error, 0))
		},
		"/opened": func(w http.ResponseWriter, r *http.Request) {
			update_key := "updated_at"
			missions := controller.DatabaseController.GetMissions(
				0,
				&update_key,
				true,
				map[string][]any{"close_at > ?": {time.Now()}},
			)
			SuccessResponse(w, missions, make([]error, 0))
		},
		"/closed": func(w http.ResponseWriter, r *http.Request) {
			update_key := "updated_at"
			missions := controller.DatabaseController.GetMissions(
				0,
				&update_key,
				true,
				map[string][]any{"close_at < ?": {time.Now()}},
			)
			SuccessResponse(w, missions, make([]error, 0))
		},
		"/available": func(w http.ResponseWriter, r *http.Request) {
			user, err := controller.DatabaseController.GetOrCreateUserFromString(r.URL.Query().Get("user"))
			if err != nil {
				CustomResponse(w, "An error occured while initializing the user", []error{err}, http.StatusInternalServerError)
				return
			}
			participations := controller.DatabaseController.GetUserParticipations(user)

			update_key := "updated_at"
			missions := controller.DatabaseController.GetMissions(
				0,
				&update_key,
				true,
				map[string][]any{"close_at > ?": {time.Now()}},
			)

			filteredMissions := []models.Mission{}
			for _, mission := range missions {
				match := false
				for _, participation := range participations {
					if participation.Mission == mission.ID {
						match = true
						break
					}
				}

				if !match {
					filteredMissions = append(filteredMissions, mission)
				}
			}
			SuccessResponse(w, filteredMissions, make([]error, 0))
		},
		"/completed": func(w http.ResponseWriter, r *http.Request) {
			user, err := controller.DatabaseController.GetOrCreateUserFromString(r.URL.Query().Get("user"))
			if err != nil {
				CustomResponse(w, "An error occured while initializing the user", []error{err}, http.StatusInternalServerError)
				return
			}
			participations := controller.DatabaseController.GetUserParticipations(user)

			update_key := "updated_at"
			missions := controller.DatabaseController.GetMissions(
				0,
				&update_key,
				true,
				map[string][]any{},
			)

			filteredMissions := []models.Mission{}
			for _, mission := range missions {
				match := false
				for _, participation := range participations {
					if participation.Mission == mission.ID {
						match = true
						break
					}
				}

				if match {
					filteredMissions = append(filteredMissions, mission)
				}
			}
			SuccessResponse(w, filteredMissions, make([]error, 0))
		},
	}
}

func init() {
	routeInitializers = append(routeInitializers, func(controller *RestController) (err error) {
		missionRouter := controller.Router.PathPrefix("/missions").Subrouter()
		missionRoutes := getMissionRoutes(controller)
		for route, handler := range missionRoutes {
			missionRouter.HandleFunc(route, handler)
		}

		return err
	})
}
