package domain

import (
	"fmt"
	"qastack-testcases/dto"
	"qastack-testcases/errs"
	"qastack-testcases/logger"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const dbTSLayout = "2006-01-02 15:04:05"

type TestRunRepositoryDb struct {
	client *sqlx.DB
}

func (tr TestRunRepositoryDb) GetTestCaseRunHistory(testCaseRunId string) ([]TestCaseRunHistory, *errs.AppError) {
	var err error
	testCaseRuns := make([]TestCaseRunHistory, 0)

	log.Info(testCaseRunId)
	getAllTestCaseRunHistory := "select status,executed_by,comments,last_execution_date  from public.test_status_records tsr where tsr.testcase_run_id =$1 order by last_execution_date desc"
	err = tr.client.Select(&testCaseRuns, getAllTestCaseRunHistory, testCaseRunId)

	if err != nil {
		fmt.Println("Error while querying test_status_records table " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return testCaseRuns, nil
}

func (tr TestRunRepositoryDb) UploadFileToS3(fileName string, project string, testRunName string, testCaseId string) *errs.AppError {
	bucketName := project + "-" + testRunName + "-" + testCaseId
	var sess = connectAWS()
	// Create AWS S3 bucket
	log.Info(project)
	log.Info(bucketName)
	// err := CreateAWSBucket(sess, bucketName)
	// if err != nil {
	// 	fmt.Println("Error while create s3 bucket " + err.Error())
	// 	return errs.NewUnexpectedError("Unexpected aws s3 event error")
	// }

	//Upload file to create S3 Bucket

	uploadErr := AddFileToS3(sess, fileName, bucketName)
	if uploadErr != nil {
		fmt.Println("Error while uploading file in s3 bucket " + uploadErr.Error())
		return errs.NewUnexpectedError("Unexpected aws s3 upload event error")
	}

	return nil

}
func (tr TestRunRepositoryDb) DownloadTestResult(project string, testRunName string, testCaseId string, fileName string) (string, *errs.AppError) {
	bucketName := project + "-" + testRunName + "-" + testCaseId
	var sess = connectAWS()
	// Create AWS S3 bucket
	log.Info(project)
	log.Info(bucketName)
	// err := CreateAWSBucket(sess, bucketName)
	// if err != nil {
	// 	fmt.Println("Error while create s3 bucket " + err.Error())
	// 	return errs.NewUnexpectedError("Unexpected aws s3 event error")
	// }

	//Upload file to create S3 Bucket

	urlStr, err := DownloadFile(sess, bucketName, fileName)
	if err != nil {
		fmt.Printf("Couldn't retrieve bucket items: %v", err)
		return "", errs.NewUnexpectedError("Unexpected aws s3 List object event error")
	}

	return urlStr, nil

}

func (tr TestRunRepositoryDb) GetTestResultsUploads(project string, testRunName string, testCaseId string) ([]dto.TestResults, *errs.AppError) {
	bucketName := project + "-" + testRunName + "-" + testCaseId
	var sess = connectAWS()
	// Create AWS S3 bucket
	log.Info(project)
	log.Info(bucketName)
	// err := CreateAWSBucket(sess, bucketName)
	// if err != nil {
	// 	fmt.Println("Error while create s3 bucket " + err.Error())
	// 	return errs.NewUnexpectedError("Unexpected aws s3 event error")
	// }

	//Upload file to create S3 Bucket
	results := []dto.TestResults{}
	s3Client := s3.New(sess)
	prefixName := ""
	bucketObjects, err := ListItems(s3Client, bucketName, prefixName)
	if err != nil {
		fmt.Printf("Couldn't retrieve bucket items: %v", err)
		return nil, errs.NewUnexpectedError("Unexpected aws s3 List object event error")
	}

	for _, item := range bucketObjects.Contents {
		fmt.Printf("Name: %s, Last Modified: %s\n", *item.Key, *item.LastModified)
		results = append(results, dto.TestResults{*item.Key})
	}

	return results, nil

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
	findAllTestTitlesTestRuns := "select ttr.id,ttr.testcase_id, t.title,ttr.status,ttr.assignee,ttr.last_execution_date,ttr.executed_by,ttr.comments  from public.testrun_testcase_records ttr join public.testcase t on t.id = ttr.testcase_id where ttr.testrun_id =$1"
	err = tr.client.Select(&testRuns, findAllTestTitlesTestRuns, id)

	if err != nil {
		fmt.Println("Error while querying testrun_testcase_records table " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return testRuns, nil
}

func (tr TestRunRepositoryDb) UpdateTestRun(id string, testRun TestRun) *errs.AppError {

	sqlTestRunTestRecordDelete := "Delete from testrun_testcase_records where testrun_id = $1 and testcase_id not in ($2)"

	res, err := tr.client.Exec(sqlTestRunTestRecordDelete, id, pq.Array(testRun.TestCases))
	if err != nil {
		return errs.NewUnexpectedError("Unexpected database error")
	}
	count, err := res.RowsAffected()
	if err != nil {
		return errs.NewUnexpectedError("Unexpected database error")
	}
	fmt.Println(count)
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

		sqlTestRunTestRecordInsert := "INSERT INTO testrun_testcase_records (testrun_id,testcase_id,status,executed_by,assignee,last_execution_date ) values ($1, $2,$3,$4,$5,$6) ON CONFLICT ON CONSTRAINT testrun_testcase_records_un DO NOTHING RETURNING id"

		_, err := tx.Exec(sqlTestRunTestRecordInsert, id, d, testRun.Status, testRun.Executed_By, testRun.Assignee, testRun.LastExecutedDate)
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

	sqlInsert := "INSERT INTO test_status_records (status, testcase_id,assignee,last_execution_date,testcase_run_id,executed_by,comments) values ($1, $2,$3,$4,$5,$6,$7) RETURNING id"

	_, err = tx.Exec(sqlInsert, testStatusRecord.Status, testStatusRecord.TestCase_Id, testStatusRecord.Assignee, testStatusRecord.LastExecutedDate, testStatusRecord.Testcase_run_Id, testStatusRecord.Executed_By, testStatusRecord.Comment)

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction into test_status_records: " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}

	update_test_status := "UPDATE testrun_testcase_records SET status = (select status from public.test_status_records tsr where testcase_run_id=$1 order by last_execution_date desc limit 1  ),executed_by =(select executed_by from public.test_status_records tsr where testcase_run_id=$2 order by last_execution_date desc limit 1),assignee =(select assignee from public.test_status_records tsr where testcase_run_id=$3 order by last_execution_date desc limit 1),comments =(select comments from public.test_status_records tsr where testcase_run_id=$2 order by last_execution_date desc limit 1) WHERE id=$4"
	_, err = tx.Exec(update_test_status, testStatusRecord.Testcase_run_Id, testStatusRecord.Testcase_run_Id, testStatusRecord.Testcase_run_Id, testStatusRecord.Testcase_run_Id)

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction into testrun_testcase_records table: " + err.Error())
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
