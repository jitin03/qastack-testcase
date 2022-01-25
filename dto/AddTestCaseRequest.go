package dto

import "github.com/jmoiron/sqlx/types"

type AddTestCaseRequest struct {
	Component_Id string         `json:"componentId"`
	Type         string         `json:"type"`
	Priority     string         `json:"priority"`
	TestStep     types.JSONText `json:"steps"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
}
