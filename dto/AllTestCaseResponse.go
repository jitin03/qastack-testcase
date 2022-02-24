package dto

import "github.com/jmoiron/sqlx/types"

type AllTestCaseResponse struct {
	TestStepsCount int            `json:"testStepCounts",omitempty`
	TotalTestCases int            `json:"total_testcases",omitempty`
	TestCaseId     string         `json:"testcaseId"`
	Type           string         `json:"type"`
	Priority       string         `json:"priority"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	TestStep       types.JSONText `json:"steps",omitempty`
	Component_id   string         `json:"component_id"`
}

type ResponseBody struct {
	Status string `json:"message"`
}
