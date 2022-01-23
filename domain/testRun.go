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
	TestRun_Id string `db:"id"`
	Name       string `db:"name"`
}

type TestRunRepository interface {
	AddTestRuns(testrun TestRun) (*TestRun, *errs.AppError)
	AllProjectTestRuns(project_id string) ([]ProjectTestRuns, *errs.AppError)
}

func (t TestRun) ToAddTestRunResponseDto() *dto.AddTestRunResponse {
	return &dto.AddTestRunResponse{t.TestRun_Id}
}
func (t ProjectTestRuns) ToProjectTestRunDto() dto.AllProjectTestRuns {
	return dto.AllProjectTestRuns{
		TestRun_Id: t.TestRun_Id,
		Name:       t.Name,
	}
}
