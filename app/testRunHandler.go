package app

import (
	"encoding/json"
	"net/http"
	"qastack-testcases/dto"
	"qastack-testcases/service"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type TestRunHandler struct {
	service service.TestRunService
}

func (t TestRunHandler) AddTestRuns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	var request dto.AddTestRunRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
	} else {

		userId, appError := t.service.AddTestRuns(request)
		if appError != nil {
			WriteResponse(w, appError.Code, appError.AsMessage())
		} else {
			WriteResponse(w, http.StatusCreated, userId)
		}
	}
}

func (t TestRunHandler) UpdateTestRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	params := mux.Vars(r)

	id := params["id"]
	var request dto.AddTestRunRequest
	type responseBody struct {
		UpdateTestRunResponse string `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
	} else {

		appError := t.service.UpdateTestRun(id, request)
		if appError != nil {
			WriteResponse(w, appError.Code, appError.AsMessage())
		} else {
			respondWithJSON(w, 200, responseBody{
				UpdateTestRunResponse: "testcase:" + id + " is updated successfully!",
			})
		}
	}
}

func (t TestRunHandler) AllProjectTestRuns(w http.ResponseWriter, r *http.Request) {
	project_id := r.URL.Query().Get("projectId")
	log.Info(project_id)
	testruns, err := t.service.AllProjectTestRuns(project_id)
	if err != nil {

		WriteResponse(w, err.Code, err.AsMessage())
	} else {

		WriteResponse(w, http.StatusOK, testruns)
	}
}

func (t TestRunHandler) GetProjectTestRun(w http.ResponseWriter, r *http.Request) {
	testRun_Id := r.URL.Query().Get("id")
	project_id := r.URL.Query().Get("projectId")
	log.Info(testRun_Id)
	testruns, err := t.service.GetProjectTestRun(project_id, testRun_Id)
	if err != nil {

		WriteResponse(w, err.Code, err.AsMessage())
	} else {

		WriteResponse(w, http.StatusOK, testruns)
	}
}

func (t TestRunHandler) GetTestCaseTitlesForTestRun(w http.ResponseWriter, r *http.Request) {
	testRun_Id := r.URL.Query().Get("id")

	log.Info(testRun_Id)
	testruns, err := t.service.GetTestCaseTitlesForTestRun(testRun_Id)
	if err != nil {

		WriteResponse(w, err.Code, err.AsMessage())
	} else {

		WriteResponse(w, http.StatusOK, testruns)
	}
}
