package service

import (
	"qastack-testcases/domain"
	"qastack-testcases/dto"
	"qastack-testcases/errs"

	log "github.com/sirupsen/logrus"
)

type DefaultTestCaseService struct {
	repo domain.TestCaseRepository
}

type TestCaseService interface {
	AddTestCase(request dto.AddTestCaseRequest, projectId string) (*dto.AddTestCaseResponse, *errs.AppError)
	ImportTestCases(request dto.ImportTestCases, projectId string) *errs.AppError
	UpdateTestCase(id string, request dto.AddTestCaseRequest) *errs.AppError
	AllTestCases(componentId string, project_id string, pageId int) ([]dto.AllTestCaseResponse, *errs.AppError)
	GetTotalTestCases(project_id string) ([]dto.AllTestCaseResponse, *errs.AppError)
	GetTestCase(testCaseId string) (*dto.AllTestCaseResponse, *errs.AppError)
}

func (s DefaultTestCaseService) AllTestCases(componentId string, project_id string, pageId int) ([]dto.AllTestCaseResponse, *errs.AppError) {

	testCases, err := s.repo.AllTestCases(componentId, project_id, pageId)
	if err != nil {
		return nil, err
	}
	response := make([]dto.AllTestCaseResponse, 0)
	for _, testCase := range testCases {
		response = append(response, testCase.ToDto())
	}
	return response, err
}

func (s DefaultTestCaseService) GetTestCase(testCaseId string) (*dto.AllTestCaseResponse, *errs.AppError) {
	testCases, err := s.repo.GetTestCase(testCaseId)
	if err != nil {
		return nil, err
	}
	response := testCases.ToDto()
	return &response, err
}

func (s DefaultTestCaseService) GetTotalTestCases(project_id string) ([]dto.AllTestCaseResponse, *errs.AppError) {

	testCases, err := s.repo.GetTotalTestCases(project_id)
	if err != nil {
		return nil, err
	}
	response := make([]dto.AllTestCaseResponse, 0)
	for _, testCase := range testCases {
		response = append(response, testCase.ToProjectTestCaseDto())
	}
	return response, err
}

func (s DefaultTestCaseService) AddTestCase(req dto.AddTestCaseRequest, projectId string) (*dto.AddTestCaseResponse, *errs.AppError) {
	log.Info("line23")
	log.Info(req.TestStep)
	testSteps := req.TestStep
	//testSteps := []dto.TestStepRequest{}
	c := domain.TestCase{

		Title:        req.Title,
		Description:  req.Description,
		Priority:     req.Priority,
		Type:         req.Type,
		Component_id: req.Component_Id,
		TestStep:     testSteps,
	}

	if newTestCase, err := s.repo.AddTestCase(c, projectId); err != nil {
		return nil, err
	} else {
		return newTestCase.ToAddTestCaseResponseDto(), nil
	}
}

func (s DefaultTestCaseService) ImportTestCases(req dto.ImportTestCases, projectId string) *errs.AppError {

	log.Info(req)
	testCases := req
	response := make([]domain.RawTestCase, 0)
	for _, testCase := range testCases {

		c := domain.RawTestCase{

			Title:          testCase.TestCaseTitle,
			Description:    testCase.Description,
			Priority:       testCase.Priority,
			Type:           testCase.Type,
			ComponentName:  testCase.Component,
			TestStep:       testCase.Steps,
			ExpectedResult: testCase.ExpectedResult,
		}

		response = append(response, c)
	}

	if err := s.repo.ImportRawTestCase(response, projectId); err != nil {
		return err
	} else {
		return nil
	}
	return nil

}

func (s DefaultTestCaseService) UpdateTestCase(id string, req dto.AddTestCaseRequest) *errs.AppError {
	testSteps := req.TestStep
	//testSteps := []dto.TestStepRequest{}
	c := domain.TestCase{

		Title:        req.Title,
		Description:  req.Description,
		Priority:     req.Priority,
		Type:         req.Type,
		Component_id: req.Component_Id,
		TestStep:     testSteps,
	}

	if err := s.repo.UpdateTestCase(id, c); err != nil {
		return err
	}
	return nil
}

func NewTestCaseService(repository domain.TestCaseRepository) DefaultTestCaseService {
	return DefaultTestCaseService{repository}
}
