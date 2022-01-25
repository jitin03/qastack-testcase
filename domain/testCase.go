package domain

import (
	"qastack-testcases/dto"
	"qastack-testcases/errs"

	"github.com/jmoiron/sqlx/types"
)

type TestCase struct {
	TestCase_Id  string `db:"id"`
	Title        string `db:"name"`
	Description  string `db:"description"`
	Component_id string `db:"component_id"`
	Type         string `db:"type"`
	Priority     string `db:"priority"`
	// TestStep     []struct {
	// 	StepDescription string `db:"stepDescription"`
	// 	ExpectedResult  string `db:"expectedResult"`
	// } `db:"steps"`
	TestStep types.JSONText `db:"steps"`
}

type OnlyTestCase struct {
	TestCase_Id  string         `db:"id"`
	Title        string         `db:"title"`
	Description  string         `db:"description"`
	Type         string         `db:"type"`
	Priority     string         `db:"priority"`
	TestStep     types.JSONText `db:"steps"`
	Component_id string         `db:"component_id"`
}

type ProjectTestCases struct {
	Total_TestCases int    `db:"total_testcases"`
	TestCase_Id     string `db:"id"`
	Title           string `db:"title"`
	Description     string `db:"description"`
	Type            string `db:"type"`
	Priority        string `db:"priority"`
}

type TestCaseRepository interface {
	AddTestCase(testcases TestCase) (*TestCase, *errs.AppError)
	UpdateTestCase(id string, testcase TestCase) *errs.AppError
	AllTestCases(componentId string, pageId int) ([]OnlyTestCase, *errs.AppError)
	GetTotalTestCases(project_id string) ([]ProjectTestCases, *errs.AppError)
	GetTestCase(testCaseId string) (*OnlyTestCase, *errs.AppError)
}

func (t TestCase) ToAddTestCaseResponseDto() *dto.AddTestCaseResponse {
	return &dto.AddTestCaseResponse{t.TestCase_Id}
}

func (t OnlyTestCase) ToDto() dto.AllTestCaseResponse {
	return dto.AllTestCaseResponse{

		TestCaseId:   t.TestCase_Id,
		Title:        t.Title,
		Description:  t.Description,
		Type:         t.Type,
		Priority:     t.Priority,
		TestStep:     t.TestStep,
		Component_id: t.Component_id,
	}
}

func (t ProjectTestCases) ToProjectTestCaseDto() dto.AllTestCaseResponse {
	return dto.AllTestCaseResponse{
		TotalTestCases: t.Total_TestCases,
		TestCaseId:     t.TestCase_Id,
		Title:          t.Title,
		Description:    t.Description,
		Type:           t.Type,
		Priority:       t.Priority,
	}
}
