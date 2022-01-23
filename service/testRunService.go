package service

import (
	"qastack-testcases/domain"
	"qastack-testcases/dto"
	"qastack-testcases/errs"

	log "github.com/sirupsen/logrus"
)

type DefaultTestRunService struct {
	repo domain.TestRunRepository
}

type TestRunService interface {
	AddTestRuns(request dto.AddTestRunRequest) (*dto.AddTestRunResponse, *errs.AppError)
	AllProjectTestRuns(project_id string) ([]dto.AllProjectTestRuns, *errs.AppError)
}

func (s DefaultTestRunService) AddTestRuns(req dto.AddTestRunRequest) (*dto.AddTestRunResponse, *errs.AppError) {
	log.Info("line23")
	log.Info(req.TestCases)

	c := domain.TestRun{
		TestRun_Id:  "",
		Name:        req.Name,
		Description: req.Description,
		Release_Id:  req.Release_Id,
		Assignee:    req.Assignee,
		TestCases:   req.TestCases,
		// TestCases: []struct {
		// 	TestCase_Id string `db:"testcase_id"`
		// }(req.TestCases),

	}

	if newTestRun, err := s.repo.AddTestRuns(c); err != nil {
		return nil, err
	} else {
		return newTestRun.ToAddTestRunResponseDto(), nil
	}
}

func (s DefaultTestRunService) AllProjectTestRuns(project_id string) ([]dto.AllProjectTestRuns, *errs.AppError) {

	testRuns, err := s.repo.AllProjectTestRuns(project_id)
	if err != nil {
		return nil, err
	}
	response := make([]dto.AllProjectTestRuns, 0)
	for _, testRun := range testRuns {
		response = append(response, testRun.ToProjectTestRunDto())
	}
	return response, err
}

func NewTestRunService(repository domain.TestRunRepository) DefaultTestRunService {
	return DefaultTestRunService{repository}
}
