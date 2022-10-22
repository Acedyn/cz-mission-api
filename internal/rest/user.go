package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getUserRoutes(controller *RestController) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/leaderboard": func(w http.ResponseWriter, r *http.Request) {
			points_key := "points"
			users := controller.DatabaseController.GetUsers(
				10,
				&points_key,
				false,
				map[string][]any{},
			)
			SuccessResponse(w, users, make([]error, 0))
		},
		"/{id}/user": func(w http.ResponseWriter, r *http.Request) {
			id := mux.Vars(r)["id"]
			user, err := controller.DatabaseController.GetOrCreateUserFromString(id)
			if err != nil {
				CustomResponse(w, "An error occured while initializing the user", []error{err}, http.StatusInternalServerError)
				return
			}
			SuccessResponse(w, *user, make([]error, 0))
		},
		"/{id}/participations": func(w http.ResponseWriter, r *http.Request) {
			id := mux.Vars(r)["id"]
			user, err := controller.DatabaseController.GetOrCreateUserFromString(id)
			if err != nil {
				CustomResponse(w, "An error occured while initializing the user", []error{err}, http.StatusInternalServerError)
				return
			}
			participations := controller.DatabaseController.GetUserParticipations(user)
			SuccessResponse(w, participations, make([]error, 0))
		},
	}
}

func init() {
	routeInitializers = append(routeInitializers, func(controller *RestController) (err error) {
		missionRouter := controller.Router.PathPrefix("/users").Subrouter()
		missionRoutes := getUserRoutes(controller)
		for route, handler := range missionRoutes {
			missionRouter.HandleFunc(route, handler)
		}

		return err
	})
}
