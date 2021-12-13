package domain

import "qastack-testcases/errs"

type TestStep struct {
	TestStep_Id    string    `db:"id"`
	Description    string `db:"stepdescription"`
	ExpectedResult string `db:"expected_result"`
	TestCase_Id    string    `db:"testcase_id"`
}

type TestStepRepository interface {
	AddTestStep(teststep TestStep) (*TestStep, *errs.AppError)
	//AllTestCases() ([]TestCase, *errs.AppError)
}
