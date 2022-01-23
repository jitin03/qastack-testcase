package domain

import (
	"fmt"
	"qastack-testcases/errs"
	"qastack-testcases/logger"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type TestCaseRepositoryDb struct {
	client *sqlx.DB
}

func (t TestCaseRepositoryDb) AddTestCase(testcases TestCase) (*TestCase, *errs.AppError) {
	// starting the database transaction block
	tx, err := t.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for bank account transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	sqlInsert := "INSERT INTO testcase (title, description,component_id,type,priority) values ($1, $2,$3,$4,$5) RETURNING id"

	_, err = tx.Exec(sqlInsert, testcases.Title, testcases.Description, testcases.Component_id, testcases.Type, testcases.Priority)

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction into test case: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// Run a query to get new test case id
	row := tx.QueryRow("SELECT id FROM testcase WHERE title=$1", testcases.Title)
	var id string
	// Store the count in the `catCount` variable
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		logger.Error("Error while getting test case id : " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	testcases.TestCase_Id = id
	log.Info(id)

	for index, d := range testcases.TestStep {
		fmt.Println(index, d)
		x := d.StepDescription
		fmt.Println(x)

		sqlTestStepInsert := "INSERT INTO teststeps (testcase_id,stepdescription, expected_result) values ($1, $2,$3) RETURNING id"

		_, err := tx.Exec(sqlTestStepInsert, id, d.StepDescription, d.ExpectedResult)
		if err != nil {
			tx.Rollback()
			logger.Error("Error while saving transaction into test case: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}

	// commit the transaction when all is good
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error while commiting transaction for testcase: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return &testcases, nil
}

func (t TestCaseRepositoryDb) AllTestCases(componentId string, pageId int) ([]OnlyTestCase, *errs.AppError) {
	var err error
	testCases := make([]OnlyTestCase, 0)
	log.Info(componentId, pageId)
	//"select id,title,description,type,priority from testcase where component_id=$1 LIMIT $2"
	findAllSql := "select count(t.testcase_id) as teststeps_count,t2.id  ,t2.title ,t2.description ,t2.type ,t2.priority from public.teststeps t  join public.testcase t2 on t.testcase_id =t2.id where t2.component_id =$1 group by t2.title ,t2.description,t2.type ,t2.priority,t2.id LIMIT $2"
	err = t.client.Select(&testCases, findAllSql, componentId, pageId)

	if err != nil {
		fmt.Println("Error while querying testcase table " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return testCases, nil
}

func (t TestCaseRepositoryDb) GetTotalTestCases(project_id string) ([]ProjectTestCases, *errs.AppError) {
	var err error
	testCases := make([]ProjectTestCases, 0)

	log.Info(project_id)
	findAllProjectTestCases := "select count(t2.id )as total_testcases ,t2.id  ,t2.title ,t2.description ,t2.type ,t2.priority  from public.testcase t2 where t2.component_id in (select id from public.component c where project_id=$1) group by t2.id"
	err = t.client.Select(&testCases, findAllProjectTestCases, project_id)

	if err != nil {
		fmt.Println("Error while querying testcase table " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return testCases, nil
}

func NewTestCaseRepositoryDb(dbClient *sqlx.DB) TestCaseRepositoryDb {
	return TestCaseRepositoryDb{dbClient}
}
