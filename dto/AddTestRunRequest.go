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

type UpdateTestStatusRequest struct {
	Testcase_run_Id string `json:"testcase_run_id"`
	Assignee        string `json:"assignee"`
	TestCase_Id     string `json:"testcase_id`
	Status          string `json:"status"`
	Executed_By     string `json:"executed_by"`
	Comments        string `json:"comments"`
}
