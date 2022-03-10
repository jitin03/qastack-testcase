package app

import (
	"encoding/json"
	"net/http"
	"qastack-testcases/dto"
	"qastack-testcases/service"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type TestCaseHandler struct {
	service service.TestCaseService
}

func (t TestCaseHandler) AddTestCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	var request dto.AddTestCaseRequest
	projectId := r.URL.Query().Get("projectId")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
	} else {

		userId, appError := t.service.AddTestCase(request, projectId)
		if appError != nil {
			WriteResponse(w, appError.Code, appError.AsMessage())
		} else {
			WriteResponse(w, http.StatusCreated, userId)
		}
	}
}
func (t TestCaseHandler) UploadTestCases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	var request dto.ImportTestCases
	projectId := r.URL.Query().Get("projectId")

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
	} else {

		appError := t.service.ImportTestCases(request, projectId)
		if appError != nil {
			WriteResponse(w, appError.Code, appError.AsMessage())
		} else {
			respondWithJSON(w, 201, dto.ResponseBody{
				Status: "Data has been inserted successfully",
			})
		}
	}
}

func (t TestCaseHandler) GetProjectTestsStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	projectId := r.URL.Query().Get("projectId")
	projectStatus, err := t.service.GetProjectTestsStatus(projectId)
	if err != nil {

		WriteResponse(w, err.Code, err.AsMessage())
	} else {

		WriteResponse(w, http.StatusOK, projectStatus)
	}

}

func (t TestCaseHandler) GetComponentTestCases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	projectId := r.URL.Query().Get("projectId")
	componentTestCases, err := t.service.GetComponentTestCases(projectId)
	if err != nil {

		WriteResponse(w, err.Code, err.AsMessage())
	} else {

		WriteResponse(w, http.StatusOK, componentTestCases)
	}

}

func (t TestCaseHandler) UpdateTestCase(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]
	var request dto.AddTestCaseRequest
	type responseBody struct {
		UpdateTestCaseResponse string `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
	} else {

		appError := t.service.UpdateTestCase(id, request)
		if appError != nil {
			WriteResponse(w, appError.Code, appError.AsMessage())
		} else {
			respondWithJSON(w, 200, responseBody{
				UpdateTestCaseResponse: "testcase:" + id + " is updated successfully!",
			})
		}
	}
}

func (t TestCaseHandler) AllTestCases(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("page")
	component := r.URL.Query().Get("componentId")
	project_id := r.URL.Query().Get("projectId")

	pageId, _ := strconv.Atoi(page)
	// componentId, _ := strconv.Atoi(component)

	components, err := t.service.AllTestCases(component, project_id, pageId)
	if err != nil {

		WriteResponse(w, err.Code, err.AsMessage())
	} else {

		WriteResponse(w, http.StatusOK, components)
	}

}

func (t TestCaseHandler) GetTestCase(w http.ResponseWriter, r *http.Request) {
	testCaseId := r.URL.Query().Get("testCaseId")

	components, err := t.service.GetTestCase(testCaseId)
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
