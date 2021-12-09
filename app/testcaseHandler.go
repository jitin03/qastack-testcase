package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"qastack-testcases/dto"
	"qastack-testcases/service"
	"strconv"
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

func (t TestCaseHandler)AllTestCases(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("page")
	component := r.URL.Query().Get("componentId")

	pageId, _ := strconv.Atoi(page)
	componentId, _ := strconv.Atoi(component)

	components, err := t.service.AllTestCases(componentId,pageId)
	if err != nil {
		fmt.Println("Inside error" + err.Message)

		WriteResponse(w, err.Code, err.AsMessage())
	} else {
		fmt.Println("Inside error")
		WriteResponse(w, http.StatusOK, components)
	}

}