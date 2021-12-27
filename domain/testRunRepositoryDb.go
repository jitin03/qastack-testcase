package domain

import (
	"fmt"
	"qastack-testcases/errs"
	"qastack-testcases/logger"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type TestRunRepositoryDb struct {
	client *sqlx.DB
}

func (tr TestRunRepositoryDb) AddTestRuns(testrun TestRun) (*TestRun, *errs.AppError) {

	// starting the database transaction block
	tx, err := tr.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for test run transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	sqlInsert := "INSERT INTO testrun (name, description,release_id,assignee) values ($1, $2,$3,$4) RETURNING id"

	_, err = tx.Exec(sqlInsert, testrun.Name, testrun.Description, testrun.Release_Id, testrun.Assignee)

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction into test case: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// Run a query to get new test case id
	row := tx.QueryRow("SELECT id FROM testrun WHERE name=$1", testrun.Name)
	var id string
	// Store the count in the `catCount` variable
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		logger.Error("Error while getting test case id : " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	testrun.TestRun_Id = id
	log.Info(id)

	for index, d := range testrun.TestCases {
		fmt.Println(index, d)
		x := d
		fmt.Println(x)

		sqlTestRunTestRecordInsert := "INSERT INTO testrun_testcase_records (testrun_id,testcase_id ) values ($1, $2) RETURNING id"

		_, err := tx.Exec(sqlTestRunTestRecordInsert, id, d)
		if err != nil {
			tx.Rollback()
			logger.Error("Error while saving transaction into test run: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}

	// commit the transaction when all is good
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error while commiting transaction for testrun: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return &testrun, nil
}

func NewTestRunRepositoryDb(dbClient *sqlx.DB) TestRunRepositoryDb {
	return TestRunRepositoryDb{dbClient}
}
