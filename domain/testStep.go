package domain

import "qastack-testcases/errs"

type TestStep struct {
	TestStep_Id    int    `db:"id"`
	Description    string `db:"stepdescription"`
	ExpectedResult string `db:"expected_result"`
	TestCase_Id    int    `db:"testcase_id"`
}

type TestStepRepository interface {
	AddTestStep(teststep TestStep) (*TestStep, *errs.AppError)
	//AllTestCases() ([]TestCase, *errs.AppError)
}
