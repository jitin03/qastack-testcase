package dto

import (
	"database/sql"

	"github.com/jmoiron/sqlx/types"
)

type AddTestCaseRequest struct {
	Component_Id string         `json:"componentId"`
	Type         string         `json:"type"`
	Priority     string         `json:"priority"`
	TestStep     types.JSONText `json:"steps"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Mode         sql.NullString `json:"mode,omitempty"`
}

type ImportTestCases []struct {
	Component      string         `json:"component"`
	TestCaseTitle  string         `json:"test_case_title"`
	Description    string         `json:"description"`
	ExpectedResult string         `json:"expected_result"`
	Steps          string         `json:"steps"`
	Mode           sql.NullString `json:"mode,omitempty"`
	Priority       string         `json:"priority"`
	Type           string         `json:"type",`
}
