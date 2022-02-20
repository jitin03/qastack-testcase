package dto

import "database/sql"

type TestCaseRunHistory struct {
	Status           string         `json:"status"`
	Executed_By      string         `json:"executed_by"`
	Comments         sql.NullString `json:"comments"`
	LastExecutedDate string         `json:"last_execution_date"`
}
