package domain

import (
	"qastack-testcases/dto"
	"qastack-testcases/errs"
)

type TestRun struct {
	TestRun_Id  string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Release_Id  string `db:"release_id"`
	Assignee    string `db:"assignee"`
	// TestCases   []struct {
	// 	TestCase_Id string `db:"testcase_id"`
	// }
	TestCases []string `db:testcases`
}

type ProjectTestRuns struct {
	TestCases  int    `db:"testcase_count"`
	Name       string `db:"name"`
	TestRun_Id string `db:"id"`
	Release_Id string `db:"release_id"`
	Assignee   string `db:"assignee"`
}

type TestRunRepository interface {
	AddTestRuns(testrun TestRun) (*TestRun, *errs.AppError)
	UpdateTestRun(id string, testrun TestRun) *errs.AppError
	AllProjectTestRuns(project_id string) ([]ProjectTestRuns, *errs.AppError)
	GetProjectTestRun(project_id string, id string) (*TestRun, *errs.AppError)
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
