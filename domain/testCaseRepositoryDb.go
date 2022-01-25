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

	sqlInsert := "INSERT INTO testcase (title, description,component_id,type,priority,steps) values ($1, $2,$3,$4,$5,$6) RETURNING id"

	var id string
	err := t.client.QueryRow(sqlInsert, testcases.Title, testcases.Description, testcases.Component_id, testcases.Type, testcases.Priority, testcases.TestStep).Scan(&id)

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		logger.Error("Error while creating new component: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}
	testcases.TestCase_Id = id

	return &testcases, nil
}

func (t TestCaseRepositoryDb) UpdateTestCase(id string, testCase TestCase) *errs.AppError {

	updateTestCaseSql := "UPDATE testcase SET title = $1 ,description = $2 ,component_id=$3, type=$4,priority=$5,steps=$6 WHERE id = $7"
	res, err := t.client.Exec(updateTestCaseSql, testCase.Title, testCase.Description, testCase.Component_id, testCase.Type, testCase.Priority, testCase.TestStep, id)
	if err != nil {
		return errs.NewUnexpectedError("Unexpected error from database")
	}
	count, err := res.RowsAffected()
	if err != nil {
		return errs.NewUnexpectedError("Unexpected error from database")
	}
	fmt.Println(count)
	return nil
}

func (t TestCaseRepositoryDb) AllTestCases(componentId string, pageId int) ([]OnlyTestCase, *errs.AppError) {
	var err error
	testCases := make([]OnlyTestCase, 0)
	log.Info(componentId, pageId)
	//"select id,title,description,type,priority from testcase where component_id=$1 LIMIT $2"
	findAllSql := "select id,title,description,type,priority,steps from public.testcase t  where component_id =$1 limit $2"
	err = t.client.Select(&testCases, findAllSql, componentId, pageId)

	if err != nil {
		fmt.Println("Error while querying testcase table " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return testCases, nil
}

func (t TestCaseRepositoryDb) GetTestCase(testCaseId string) (*OnlyTestCase, *errs.AppError) {
	var err error
	var testCases OnlyTestCase
	log.Info(testCaseId)

	//"select id,title,description,type,priority from testcase where component_id=$1 LIMIT $2"
	findAllSql := "select id,title,description,type,priority,steps,component_id from public.testcase t  where id =$1"
	err = t.client.Get(&testCases, findAllSql, testCaseId)

	if err != nil {
		fmt.Println("Error while querying testcase table " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return &testCases, nil
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
