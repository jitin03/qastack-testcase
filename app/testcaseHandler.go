package app

import (
	"encoding/json"
	"net/http"
	"qastack-testcases/dto"
	"qastack-testcases/service"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type TestCaseHandler struct {
	service service.TestCaseService
}

func (t TestCaseHandler) AddTestCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	var request dto.AddTestCaseRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
	} else {

		userId, appError := t.service.AddTestCase(request)
		if appError != nil {
			WriteResponse(w, appError.Code, appError.AsMessage())
		} else {
			WriteResponse(w, http.StatusCreated, userId)
		}
	}
}

func (t TestCaseHandler) AllTestCases(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("page")
	component := r.URL.Query().Get("componentId")

	pageId, _ := strconv.Atoi(page)
	// componentId, _ := strconv.Atoi(component)

	components, err := t.service.AllTestCases(component, pageId)
	if err != nil {

		WriteResponse(w, err.Code, err.AsMessage())
	} else {

		WriteResponse(w, http.StatusOK, components)
	}

}

func (t TestCaseHandler) GetTotalTestCases(w http.ResponseWriter, r *http.Request) {
	project_id := r.URL.Query().Get("projectId")
	log.Info(project_id)
	components, err := t.service.GetTotalTestCases(project_id)
	if err != nil {

		WriteResponse(w, err.Code, err.AsMessage())
	} else {

		WriteResponse(w, http.StatusOK, components)
	}
}
