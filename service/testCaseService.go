package service

import (
	log "github.com/sirupsen/logrus"
	"qastack-testcases/domain"
	"qastack-testcases/dto"
	"qastack-testcases/errs"
)

type DefaultTestCaseService struct {
	repo domain.TestCaseRepository
}



type TestCaseService interface {
	AddTestCase(request dto.AddTestCaseRequest) (*dto.AddTestCaseResponse, *errs.AppError)
	AllTestCases(componentId int,pageId int) ([]dto.AllTestCaseResponse, *errs.AppError)
}


func (s DefaultTestCaseService) AllTestCases(componentId int , pageId int) ([]dto.AllTestCaseResponse, *errs.AppError) {

	testCases, err := s.repo.AllTestCases(componentId,pageId)
	if err != nil {
		return nil, err
	}
	response := make([]dto.AllTestCaseResponse, 0)
	for _, testCase := range testCases {
		response = append(response, testCase.ToDto())
	}
	return response, err
}

func (s DefaultTestCaseService) AddTestCase(req dto.AddTestCaseRequest) (*dto.AddTestCaseResponse, *errs.AppError) {
	log.Info("line23")
	log.Info(req.TestStep)
	//testSteps := []dto.TestStepRequest{}
	c := domain.TestCase{

		Title:       req.Title,
		Description: req.Description,
		Priority: req.Priority,
		Type: req.Type,
		Component_id: req.Component_Id,
		TestStep: []struct {
			StepDescription string `db:"stepDescription"`
			ExpectedResult  string `db:"expectedResult"`
		}(req.TestStep),
		//components := make([]Component, 0)
		//TestStep : append([]domain.TestStep,req.TestStep...),

	}


	if newTestCase, err := s.repo.AddTestCase(c); err != nil {
		return nil, err
	} else {
		return newTestCase.ToAddTestCaseResponseDto(), nil
	}
}

func NewTestCaseService(repository domain.TestCaseRepository) DefaultTestCaseService {
	return DefaultTestCaseService{repository}
}