package domain

import (
	"database/sql"
	"qastack-testcases/dto"
	"qastack-testcases/errs"
)

type TestRun struct {
	TestRun_Id       string   `db:"id"`
	Name             string   `db:"name"`
	Description      string   `db:"description"`
	Release_Id       string   `db:"release_id"`
	Assignee         string   `db:"assignee"`
	Status           string   `db:"status"`
	Executed_By      string   `db:"executed_by"`
	LastExecutedDate string   `db:"last_execution_date"`
	TestCases        []string `db:testcases`
}

type TestStatusRecord struct {
	TestStatusId     string         `db:"id"`
	Testcase_run_Id  string         `db:"testcase_run_id"`
	Assignee         string         `db:"assignee"`
	TestCase_Id      string         `db:"testcase_id`
	Status           string         `db:"status"`
	Executed_By      string         `db:"executed_by"`
	LastExecutedDate string         `db:"last_execution_date"`
	Comment          sql.NullString `db:"comments"`
}

type ProjectTestRuns struct {
	TestCases  int    `db:"testcase_count"`
	Name       string `db:"name"`
	TestRun_Id string `db:"id"`
	Release_Id string `db:"release_id"`
	Assignee   string `db:"assignee"`
}

type TestCaseRunHistory struct {
	Status           string         `db:"status"`
	Executed_By      string         `db:"executed_by"`
	Comments         sql.NullString `db:"comments"`
	LastExecutedDate string         `db:"last_execution_date"`
}

type TestCaseTitleTestRuns struct {
	TestCaseRunId    string         `db:"id"`
	TestCase_Id      string         `db:"testcase_id"`
	TestCase_Title   string         `db:"title"`
	Status           string         `db:"status"`
	Assignee         string         `db:"assignee"`
	LastExecutedDate string         `db:"last_execution_date"`
	Executed_By      string         `db:"executed_by"`
	Comment          sql.NullString `db:"comments"`
}

type TestRunRepository interface {
	AddTestRuns(testrun TestRun) (*TestRun, *errs.AppError)
	UpdateTestStatus(testStatus TestStatusRecord, testRunId string) *errs.AppError
	UpdateTestRun(id string, testrun TestRun) *errs.AppError
	AllProjectTestRuns(project_id string) ([]ProjectTestRuns, *errs.AppError)
	GetTestCaseTitlesForTestRun(id string) ([]TestCaseTitleTestRuns, *errs.AppError)
	GetProjectTestRun(project_id string, id string) (*TestRun, *errs.AppError)
	GetTestCaseRunHistory(testCaseRunId string) ([]TestCaseRunHistory, *errs.AppError)
	UploadFileToS3(fileName string, project string, testRunName string, testCaseId string) *errs.AppError
	GetTestResultsUploads(project string, testRunName string, testCaseId string) ([]dto.TestResults, *errs.AppError)
	DownloadTestResult(project string, testRunName string, testCaseId string, fileName string) (string, *errs.AppError)
}

func (t TestRun) ToAddTestRunResponseDto() *dto.AddTestRunResponse {
	return &dto.AddTestRunResponse{t.TestRun_Id}
}
func (t ProjectTestRuns) ToProjectTestRunDto() dto.AllProjectTestRuns {
	return dto.AllProjectTestRuns{
		TestRun_Id: t.TestRun_Id,
		Name:       t.Name,
		Release_Id: t.Release_Id,
		TestCases:  t.TestCases,
		Assignee:   t.Assignee,
	}
}

func (t TestCaseRunHistory) ToTestCaseRunHistory() dto.TestCaseRunHistory {
	return dto.TestCaseRunHistory{
		Status:           t.Status,
		Executed_By:      t.Executed_By,
		Comments:         t.Comments,
		LastExecutedDate: t.LastExecutedDate,
	}
}

func (t TestCaseTitleTestRuns) ToTestCaseTitleTestRunDto() dto.GetTestCaseTitleTestRun {
	return dto.GetTestCaseTitleTestRun{
		TestCaseRunId:    t.TestCaseRunId,
		TestCase_Id:      t.TestCase_Id,
		TestCase_Title:   t.TestCase_Title,
		Status:           t.Status,
		Assignee:         t.Assignee,
		LastExecutedDate: t.LastExecutedDate,
		Executed_By:      t.Executed_By,
		Comment:          t.Comment,
	}
}

type TestResults struct {
	FileName string `db:"file_name"`
}

func (t TestResults) ToTestResults() dto.TestResults {
	return dto.TestResults{
		FileName: t.FileName,
	}
}

func (t TestRun) ToTestRunDto() dto.GetTestRun {
	return dto.GetTestRun{
		TestRun_Id:  t.TestRun_Id,
		Name:        t.Name,
		Release_Id:  t.Release_Id,
		TestCases:   t.TestCases,
		Assignee:    t.Assignee,
		Description: t.Description,
	}
}
