package dto

type AddTestCaseResponse struct {
	TestCase_Id string `json:"testcase_id"`
}

type ProjectStatus struct {
	Count  string `json:"test_case_count"`
	Status string `json:"status"`
}

type ComponentTestCases struct {
	Count         string `json:"test_case_count"`
	ComponentName string `json:"component_name"`
}

type TestcaseProgressResponse struct {
	CreateDate    string `json:"date"`
	TestCaseCount string `json:"testcasecount"`
}
type AddTestRunResponse struct {
	TestRun_Id string `json:"testrun_id"`
}
