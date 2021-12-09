package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"net/http"
	"qastack-testcases/domain"
	"qastack-testcases/service"
)

func getDbClient() *sqlx.DB {
	client, err := sqlx.ConnectContext(context.Background(), "postgres", "host=localhost port=5433 user=postgres dbname=postgres sslmode=disable password=qastack")
	if err != nil {
		panic(err)
	}
	return client
}


func Start() {

	//sanityCheck()

	router := mux.NewRouter()
	dbClient := getDbClient()

	router.Use()
	testcaseRepositoryDb := domain.NewTestCaseRepositoryDb(dbClient)
	////wiring
	////u := ComponentHandler{service.NewUserService(userRepositoryDb,domain.GetRolePermissions())}
	//
	t := TestCaseHandler{service.NewTestCaseService(testcaseRepositoryDb)}
	//
	// define routes

	router.HandleFunc("/api/testcases/health", func (w http.ResponseWriter,r *http.Request) {
		json.NewEncoder(w).Encode("Running...")
	})
	//
	router.
		HandleFunc("/api/testcase/add", t.AddTestCase).
		Methods(http.MethodPost).Name("AddTestCase")

	router.
		HandleFunc("/api/testcases",t.AllTestCases).
		Methods(http.MethodGet).Name("AllTestCases")
	//
	//router.
	//	HandleFunc("/api/component/delete/{id}", c.DeleteComponent).
	//	Methods(http.MethodDelete).Name("DeleteComponent")
	//
	//router.
	//	HandleFunc("/api/component/update/{id}", c.UpdateComponent).
	//	Methods(http.MethodPut).Name("UpdateComponent")

	cor := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedHeaders: []string{ "Content-Type", "Authorization","Referer"},
		AllowCredentials: true,
		AllowedMethods: []string{"GET","PUT","DELETE","POST"},
	})

	handler := cor.Handler(router)
	am := AuthMiddleware{domain.NewAuthRepository()}
	router.Use(am.authorizationHandler())

	//logger.Info(fmt.Sprintf("Starting server on %s:%s ...", address, port))
	if err := http.ListenAndServe(":8092", handler); err != nil {
		fmt.Println("Failed to set up server")

	}

}