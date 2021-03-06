package domain

import (
	"qastack-testcases/dto"
	"qastack-testcases/errs"

	"github.com/jmoiron/sqlx/types"
)

type Tenant struct {
	TenantId string `db:"tenant_id"`
	Config   string `db:"config"`
}

type TestCase struct {
	TestCase_Id  string `db:"id"`
	Title        string `db:"name"`
	Description  string `db:"description"`
	Component_id string `db:"component_id"`
	Type         string `db:"type"`
	Priority     string `db:"priority"`
	Mode         string `db:"mode"`
	CreateDate   string `db:"create_date"`
	UpdateDate   string `db:"update_date"`
	// TestStep     []struct {
	// 	StepDescription string `db:"stepDescription"`
	// 	ExpectedResult  string `db:"expectedResult"`
	// } `db:"steps"`
	TestStep types.JSONText `db:"steps"`
}

type Steps struct {
	ExpectedResult  string `json:"expectedResult"`
	StepDescription string `json:"stepDescription"`
}

type RawTestCase struct {
	TestCase_Id   string `db:"id"`
	Title         string `db:"name"`
	Description   string `db:"description"`
	ComponentName string `db:"ComponentName"`
	CreateDate    string `db:"create_date"`
	UpdateDate    string `db:"update_date"`
	Type          string `db:"type"`
	Priority      string `db:"priority"`
	Mode          string `db:"mode"`
	// TestStep     []struct {
	// 	StepDescription string `db:"stepDescription"`
	// 	ExpectedResult  string `db:"expectedResult"`
	// } `db:"steps"`
	TestStep       string `db:"steps"`
	ExpectedResult string `db:"expectedResult"`
}

type OnlyTestCase struct {
	TestCase_Id  string         `db:"id"`
	Title        string         `db:"title"`
	Description  string         `db:"description"`
	Type         string         `db:"type"`
	Priority     string         `db:"priority"`
	TestStep     types.JSONText `db:"steps"`
	Component_id string         `db:"component_id"`
	Mode         string         `db:"mode"`
}

type ProjectStatus struct {
	Count  string `db:"count"`
	Status string `db:"status"`
}

type TestcaseProgress struct {
	CreateDate    string `db:"date"`
	TestCaseCount string `db:"testcasecount"`
}

type ComponentTestCases struct {
	Count         string `db:"count"`
	ComponentName string `db:"name"`
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
	AddTestCase(testcases TestCase, projectId string) (*TestCase, *errs.AppError)
	ImportRawTestCase(testcases []RawTestCase, projectId string) *errs.AppError
	UpdateTestCase(id string, testcase TestCase) *errs.AppError
	AllTestCases(componentId string, project_id string, pageId int) ([]OnlyTestCase, *errs.AppError)
	GetProjectTestsStatus(projectId string) ([]ProjectStatus, *errs.AppError)
	GetProjectTestsProgress(projectId string) ([]TestcaseProgress, *errs.AppError)
	GetComponentTestCases(projectId string) ([]ComponentTestCases, *errs.AppError)
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
		Mode:         t.Mode,
	}
}

func (t ProjectStatus) ToProjectStatusDto() dto.ProjectStatus {
	return dto.ProjectStatus{

		Count:  t.Count,
		Status: t.Status,
	}
}

func (t TestcaseProgress) ToProjectTestsProgressDto() dto.TestcaseProgressResponse {
	return dto.TestcaseProgressResponse{

		CreateDate:    t.CreateDate,
		TestCaseCount: t.TestCaseCount,
	}
}

func (t ComponentTestCases) ToComponentTestCaseDto() dto.ComponentTestCases {
	return dto.ComponentTestCases{

		Count:         t.Count,
		ComponentName: t.ComponentName,
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
