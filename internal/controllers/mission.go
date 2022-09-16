package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cardboard-citizens/cz-goodboard-api/internal/models"
	"github.com/cardboard-citizens/cz-goodboard-api/internal/utils"
)

func (server *Server) CreateMission(w http.ResponseWriter, r *http.Request) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.JsonError(w, http.StatusUnprocessableEntity, err)
		return err
	}

	mission := models.Mission{}
	err = json.Unmarshal(body, &mission)
	if err != nil {
		utils.JsonError(w, http.StatusUnprocessableEntity, err)
		return err
	}

	err = mission.Setup()
	if err != nil {
		utils.JsonError(w, http.StatusUnprocessableEntity, err)
		return err
	}

	err = mission.Create(server.DB)
	if err != nil {
		utils.JsonError(w, http.StatusUnprocessableEntity, err)
		return err
	}

	utils.JsonResponse(w, http.StatusCreated, mission)
	return err
}
