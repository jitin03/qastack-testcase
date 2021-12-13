package domain

import (
	"qastack-testcases/dto"
	"qastack-testcases/errs"
)

type TestCase struct {
	TestCase_Id  string      `db:"id"`
	Title        string   `db:"name"`
	Description  string   `db:"description"`
	Component_id string      `db:"component_id"`
	Type         string   `db:"type"`
	Priority     string   `db:"priority"`
	TestStep     []struct {
		StepDescription string `db:"stepDescription"`
		ExpectedResult  string `db:"expectedResult"`
	} `db:"steps"'`
}

type OnlyTestCase struct {
	TestSteps    int	`db:"teststeps_count"`
	TestCase_Id  string      `db:"id"`
	Title        string   `db:"title"`
	Description  string   `db:"description"`
	Type         string   `db:"type"`
	Priority     string   `db:"priority"`
}

type TestCaseRepository interface {
	AddTestCase(testcases TestCase) (*TestCase, *errs.AppError)
	AllTestCases(componentId string, pageId int) ([]OnlyTestCase, *errs.AppError)
}

func (t TestCase) ToAddTestCaseResponseDto() *dto.AddTestCaseResponse {
	return &dto.AddTestCaseResponse{t.TestCase_Id}
}

func (t OnlyTestCase) ToDto() dto.AllTestCaseResponse {
	return dto.AllTestCaseResponse{
		TestStepsCount: t.TestSteps,
		TestCaseId: t.TestCase_Id,
		Title:       t.Title,
		Description: t.Description,
		Type: t.Type,
		Priority: t.Priority,
	}
}