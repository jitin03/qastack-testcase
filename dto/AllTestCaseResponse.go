package dto

type AllTestCaseResponse struct {
	TestStepsCount int    `json:"testStepCounts"`
	TestCaseId     int    `json:"testcaseId"`
	Type           string `json:"type"`
	Priority       string `json:"priority"`
	Title          string `json:"title"`
	Description    string `json:"description"`
}
