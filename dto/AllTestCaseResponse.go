package dto

type AllTestCaseResponse struct {
	TestStepsCount int    `json:"testStepCounts",omitempty`
	TotalTestCases int    `json:"total_testcases",omitempty`
	TestCaseId     string `json:"testcaseId"`
	Type           string `json:"type"`
	Priority       string `json:"priority"`
	Title          string `json:"title"`
	Description    string `json:"description"`
}
