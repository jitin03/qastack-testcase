package dto

type AddTestCaseRequest struct {
	Component_Id string `json:"componentId"`
	Type         string `json:"type"`
	Priority     string `json:"priority"`
	TestStep     []struct {
		StepDescription string `json:"stepDescription"`
		ExpectedResult  string `json:"expectedResult"`
	} `json:"Steps"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
