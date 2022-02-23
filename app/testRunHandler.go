package app

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
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

func (t TestRunHandler) GetTestCaseRunHistory(w http.ResponseWriter, r *http.Request) {

	testCaseRunId := r.URL.Query().Get("testCaseRunId")

	testruns, err := t.service.GetTestCaseRunHistory(testCaseRunId)
	if err != nil {

		WriteResponse(w, err.Code, err.AsMessage())
	} else {

		WriteResponse(w, http.StatusOK, testruns)
	}

}

func (t TestRunHandler) UpdateTestStatus(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	testRunId := params["testRunId"]
	type responseBody struct {
		UpdateTestStatus string `json:"status"`
	}

	var request dto.UpdateTestStatusRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
	} else {

		appError := t.service.UpdateTestStatusRequest(request, testRunId)
		if appError != nil {
			WriteResponse(w, appError.Code, appError.AsMessage())
		} else {
			respondWithJSON(w, 200, responseBody{
				UpdateTestStatus: "status has updated successfully!",
			})
		}
	}
}

func (t TestRunHandler) UploadResult(w http.ResponseWriter, r *http.Request) {

	// Parse request body as multipart form data with 32MB max memory
	r.ParseMultipartForm(32 << 20)
	project := r.URL.Query().Get("projectId")
	testRunName := r.URL.Query().Get("testRunId")
	testCaseId := r.URL.Query().Get("testCaseId")

	// Get file from Form
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Error(err)
		WriteResponse(w, http.StatusInternalServerError, ".errors.upload_image.cannot_read_file")
		return
	}
	defer file.Close()

	// Create file locally
	dst, err := os.Create(handler.Filename)
	if err != nil {
		log.Error(err)
		WriteResponse(w, http.StatusInternalServerError, ".errors.upload_image.cannot_create_local_file")
		return
	}
	defer dst.Close()

	// Copy the uploaded file data to the newly created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		log.Error(err)
		WriteResponse(w, http.StatusInternalServerError, ".errors.upload_image.cannot_copy_to_file")
		return
	}

	//Create session with AWS

	//Create a bucket in AWS S3 BUCKET

	uploadFileError := t.service.UploadFileToS3(handler.Filename, project, testRunName, testCaseId)

	if uploadFileError != nil {
		log.Error(err)
		WriteResponse(w, http.StatusInternalServerError, "File Upload error")
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (t TestRunHandler) DownloadTestResult(w http.ResponseWriter, r *http.Request) {
	project := r.URL.Query().Get("projectId")
	fileName := r.URL.Query().Get("fileName")

	testRunName := r.URL.Query().Get("testRunId")
	testCaseId := r.URL.Query().Get("testCaseId")

	result, uploadFileError := t.service.DownloadTestResult(project, testRunName, testCaseId, fileName)

	if uploadFileError != nil {
		log.Error(uploadFileError)
		WriteResponse(w, http.StatusInternalServerError, "File Upload error")
		return
	}
	type responseBody struct {
		Url string `json:"url"`
	}

	respondWithJSON(w, 200, responseBody{
		Url: result,
	})

}

func (t TestRunHandler) GetTestResultsUploads(w http.ResponseWriter, r *http.Request) {

	project := r.URL.Query().Get("projectId")
	testRunName := r.URL.Query().Get("testRunId")
	testCaseId := r.URL.Query().Get("testCaseId")

	result, uploadFileError := t.service.GetTestResultsUploads(project, testRunName, testCaseId)

	if uploadFileError != nil {
		log.Error(uploadFileError)
		WriteResponse(w, http.StatusInternalServerError, "File Upload error")
		return
	}
	WriteResponse(w, http.StatusOK, result)

}
