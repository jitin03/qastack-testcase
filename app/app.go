package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"qastack-testcases/domain"
	"qastack-testcases/logger"
	"qastack-testcases/service"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func getDbClient() *sqlx.DB {

	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbAddr := os.Getenv("DB_ADDR")
	//dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	os.Setenv("TENANT_NAME", "postgres")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbAddr, 5432, dbUser, dbPasswd, dbName)
	logger.Info(psqlInfo)
	client, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	//client, err := sqlx.ConnectContext(context.Background(), "postgres",os.Getenv("DATABASE_URL") )
	//if err != nil {
	//	panic(err)
	//}
	return client
}

func Start() {

	//sanityCheck()

	router := mux.NewRouter()
	dbClient := getDbClient()

	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Referer"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "PUT", "DELETE", "POST"},
	})

	handler := cor.Handler(router)
	am := AuthMiddleware{domain.NewAuthRepository()}
	router.Use(am.authorizationHandler())

	router.Use()

	testcaseRepositoryDb := domain.NewTestCaseRepositoryDb(dbClient)
	testRunRepositoryDb := domain.NewTestRunRepositoryDb(dbClient)
	////wiring
	////u := ComponentHandler{service.NewUserService(userRepositoryDb,domain.GetRolePermissions())}
	//
	t := TestCaseHandler{service.NewTestCaseService(testcaseRepositoryDb)}

	tr := TestRunHandler{service.NewTestRunService(testRunRepositoryDb)}
	//
	// define routes

	router.HandleFunc("/api/testcases/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Running...")
	})
	//
	router.
		HandleFunc("/api/testcase/add", t.AddTestCase).
		Methods(http.MethodPost).Name("AddTestCase")

	router.
		HandleFunc("/api/testcases", t.AllTestCases).
		Methods(http.MethodGet).Name("AllTestCases")

	router.
		HandleFunc("/api/testcase", t.GetTestCase).
		Methods(http.MethodGet).Name("GetTestCase")

	router.
		HandleFunc("/api/testruns/project", tr.AllProjectTestRuns).
		Methods(http.MethodGet).Name("AllProjectTestRuns")

	router.
		HandleFunc("/api/testrun/project", tr.GetProjectTestRun).
		Methods(http.MethodGet).Name("GetProjectTestRun")

	router.
		HandleFunc("/api/testcases/project", t.GetTotalTestCases).
		Methods(http.MethodGet).Name("GetTotalTestCases")

	router.
		HandleFunc("/api/testcase/update/{id}", t.UpdateTestCase).
		Methods(http.MethodPut).Name("UpdateTestCase")

	router.HandleFunc("/api/testrun/add", tr.AddTestRuns).Methods(http.MethodPost).Name("AddTestRuns")

	router.
		HandleFunc("/api/testrun/update/{id}", tr.UpdateTestRun).
		Methods(http.MethodPut).Name("UpdateTestRun")
	router.
		HandleFunc("/api/testrun/testcases", tr.GetTestCaseTitlesForTestRun).
		Methods(http.MethodGet).Name("GetTestCaseTitlesForTestRun")

	router.
		HandleFunc("/api/testrun/{testRunId}/test/update/status", tr.UpdateTestStatus).
		Methods(http.MethodPost).Name("UpdateTestStatus")

	router.
		HandleFunc("/api/testrun/testcase-history", tr.GetTestCaseRunHistory).
		Methods(http.MethodGet).Name("GetTestCaseRunHistory")

	router.
		HandleFunc("/api/testrun/result/upload", tr.UploadResult).
		Methods(http.MethodPost).Name("UploadResult")

	router.
		HandleFunc("/api/testrun/results/download", tr.GetTestResultsUploads).
		Methods(http.MethodGet).Name("GetTestResultsUploads")

	router.
		HandleFunc("/api/testrun/result/download", tr.DownloadTestResult).
		Methods(http.MethodGet).Name("DownloadTestResult")
	router.
		HandleFunc("/api/testcase/upload", t.UploadTestCases).
		Methods(http.MethodPost).Name("UploadTestCases")

	router.
		HandleFunc("/api/testcase/project/status", t.GetProjectTestsStatus).
		Methods(http.MethodGet).Name("GetProjectTestsStatus")

	router.
		HandleFunc("/api/testcases/components", t.GetComponentTestCases).
		Methods(http.MethodGet).Name("GetComponentTestCases")

	//logger.Info(fmt.Sprintf("Starting server on %s:%s ...", address, port))
	if err := http.ListenAndServe(":8092", handler); err != nil {
		fmt.Println("Failed to set up server")

	}

}
