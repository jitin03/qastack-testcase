package dto

type AddTestRunRequest struct {
	Release_Id string   `json:"release_id"`
	Assignee   string   `json:"assignee"`
	TestCases  []string `json:"testcases`
	// TestCases  []struct {
	// 	TestCase_Id string
	// } `json:"testcases"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
