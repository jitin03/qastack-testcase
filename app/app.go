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

	router.HandleFunc("/api/testrun/add", tr.AddTestRuns).Methods(http.MethodPost).Name("AddTestRuns")

	router.
		HandleFunc("/api/testcases", t.AllTestCases).
		Methods(http.MethodGet).Name("AllTestCases")

	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Referer"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "PUT", "DELETE", "POST"},
	})

	handler := cor.Handler(router)
	am := AuthMiddleware{domain.NewAuthRepository()}
	router.Use(am.authorizationHandler())

	//logger.Info(fmt.Sprintf("Starting server on %s:%s ...", address, port))
	if err := http.ListenAndServe(":8092", handler); err != nil {
		fmt.Println("Failed to set up server")

	}

}
