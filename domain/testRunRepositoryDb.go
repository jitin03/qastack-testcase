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

func (tr TestRunRepositoryDb) AllProjectTestRuns(project_id string) ([]ProjectTestRuns, *errs.AppError) {
	var err error
	testRuns := make([]ProjectTestRuns, 0)

	log.Info(project_id)
	findAllProjectTestRuns := "select count(ttr.testcase_id) as testcase_count,  name ,release_id,t.id as id ,t.assignee  from public.testrun t   join public.testrun_testcase_records ttr   on t.id =ttr.testrun_id where t.release_id in (select id from public.releases r  where project_id=$1) group by t.name ,t.release_id,ttr.testrun_id ,t.id "
	err = tr.client.Select(&testRuns, findAllProjectTestRuns, project_id)

	if err != nil {
		fmt.Println("Error while querying testcase table " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return testRuns, nil
}

func (tr TestRunRepositoryDb) GetProjectTestRun(project_id string, id string) (*TestRun, *errs.AppError) {
	var err error
	var testRun TestRun

	log.Info(project_id)
	findProjectTestRun := "select   t.name,t.description ,t.release_id,t.id ,t.assignee  from public.testrun t   join public.testrun_testcase_records ttr   on t.id =ttr.testrun_id where t.release_id in (select id from public.releases r  where project_id=$1) and t.id =$2 group by t.name ,t.release_id,ttr.testrun_id ,t.id"
	err = tr.client.Get(&testRun, findProjectTestRun, project_id, id)

	if err != nil {
		fmt.Println("Error while querying testrun table " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	} else {

		findTestCasesForTestRun := "select ttr.testcase_id  from public.testrun_testcase_records ttr where ttr.testrun_id =$1"
		err = tr.client.Select(&testRun.TestCases, findTestCasesForTestRun, id)

		if err != nil {
			fmt.Println("Error while querying testrun_testcase_records table " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}

	return &testRun, nil
}

func (tr TestRunRepositoryDb) GetTestCaseTitlesForTestRun(id string) ([]TestCaseTitleTestRuns, *errs.AppError) {

	var err error
	testRuns := make([]TestCaseTitleTestRuns, 0)

	log.Info(id)
	findAllTestTitlesTestRuns := "select ttr.id,ttr.testcase_id, t.title,ttr.status,ttr.assignee,ttr.last_execution_date,ttr.executed_by  from public.testrun_testcase_records ttr join public.testcase t on t.id = ttr.testcase_id where ttr.testrun_id =$1"
	err = tr.client.Select(&testRuns, findAllTestTitlesTestRuns, id)

	if err != nil {
		fmt.Println("Error while querying testrun_testcase_records table " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return testRuns, nil
}

func (tr TestRunRepositoryDb) UpdateTestRun(id string, testRun TestRun) *errs.AppError {

	// starting the database transaction block
	tx, err := tr.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for test run transaction: " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}

	updateTestRunSql := "UPDATE testrun SET name = $1 ,description = $2 ,release_id=$3, assignee=$4 WHERE id = $5"
	_, err = tx.Exec(updateTestRunSql, testRun.Name, testRun.Description, testRun.Release_Id, testRun.Assignee, id)
	if err != nil {
		return errs.NewUnexpectedError("Unexpected error from database")
	}

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction into test case: " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}

	for index, d := range testRun.TestCases {
		fmt.Println(index, d)
		x := d
		fmt.Println(x)

		sqlTestRunTestRecordInsert := "INSERT INTO testrun_testcase_records (testrun_id,testcase_id ) values ($1, $2) RETURNING id"

		_, err := tx.Exec(sqlTestRunTestRecordInsert, id, d)
		if err != nil {
			tx.Rollback()
			logger.Error("Error while saving transaction into test run: " + err.Error())
			return errs.NewUnexpectedError("Unexpected database error")
		}
	}

	// commit the transaction when all is good
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error while commiting transaction for testrun: " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}

	return nil
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

		sqlTestRunTestRecordInsert := "INSERT INTO testrun_testcase_records (testrun_id,testcase_id,status,executed_by,last_execution_date,assignee ) values ($1, $2,$3,$4,$5,$6) RETURNING id"

		_, err := tx.Exec(sqlTestRunTestRecordInsert, id, d, testrun.Status, testrun.Executed_By, testrun.LastExecutedDate, testrun.Assignee)
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
func (tr TestRunRepositoryDb) UpdateTestStatus(testStatusRecord TestStatusRecord, testRunId string) *errs.AppError {

	// starting the database transaction block
	tx, err := tr.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for test status transaction: " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}

	sqlInsert := "INSERT INTO test_status_records (status, testcase_id,assignee,last_execution_date,testcase_run_id,executed_by) values ($1, $2,$3,$4,$5,$6) RETURNING id"

	_, err = tx.Exec(sqlInsert, testStatusRecord.Status, testStatusRecord.TestCase_Id, testStatusRecord.Assignee, testStatusRecord.LastExecutedDate, testStatusRecord.Testcase_run_Id, testStatusRecord.Executed_By)

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction into test case: " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}

	// commit the transaction when all is good
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error while commiting transaction for testrun: " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}
	return nil
}
func NewTestRunRepositoryDb(dbClient *sqlx.DB) TestRunRepositoryDb {
	return TestRunRepositoryDb{dbClient}
}
