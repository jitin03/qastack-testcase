package dto

type AllProjectTestRuns struct {
	TestCases  int    `json:"testcase_count"`
	Name       string `json:"name"`
	TestRun_Id string `json:"id"`
	Release_Id string `json:"release_id"`
	Assignee   string `json:"assignee"`
}

type GetTestRun struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	TestRun_Id  string   `json:"id"`
	Release_Id  string   `json:"release_id"`
	Assignee    string   `json:"assignee"`
	TestCases   []string `json:"testcases"`
}
