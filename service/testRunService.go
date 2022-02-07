package service

import (
	"qastack-testcases/domain"
	"qastack-testcases/dto"
	"qastack-testcases/errs"
	"time"

	log "github.com/sirupsen/logrus"
)

const dbTSLayout = "2006-01-02 15:04:05"

type DefaultTestRunService struct {
	repo domain.TestRunRepository
}

type TestRunService interface {
	AddTestRuns(request dto.AddTestRunRequest) (*dto.AddTestRunResponse, *errs.AppError)
	UpdateTestRun(id string, request dto.AddTestRunRequest) *errs.AppError
	AllProjectTestRuns(project_id string) ([]dto.AllProjectTestRuns, *errs.AppError)
	GetProjectTestRun(project_id string, id string) (*dto.GetTestRun, *errs.AppError)
	GetTestCaseTitlesForTestRun(id string) ([]dto.GetTestCaseTitleTestRun, *errs.AppError)
	UpdateTestStatusRequest(request dto.UpdateTestStatusRequest, testRunId string) *errs.AppError
}

func (s DefaultTestRunService) AddTestRuns(req dto.AddTestRunRequest) (*dto.AddTestRunResponse, *errs.AppError) {
	log.Info("line23")
	log.Info(req.TestCases)

	c := domain.TestRun{
		TestRun_Id:       "",
		Name:             req.Name,
		Description:      req.Description,
		Release_Id:       req.Release_Id,
		Assignee:         req.Assignee,
		TestCases:        req.TestCases,
		Status:           "Unexecuted",
		Executed_By:      req.Assignee,
		LastExecutedDate: time.Now().Format(dbTSLayout),
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

func (s DefaultTestRunService) UpdateTestStatusRequest(req dto.UpdateTestStatusRequest, testRunId string) *errs.AppError {

	status := domain.TestStatusRecord{
		TestStatusId:     "",
		Testcase_run_Id:  req.Testcase_run_Id,
		Assignee:         req.Assignee,
		TestCase_Id:      req.TestCase_Id,
		Status:           req.Status,
		Executed_By:      req.Executed_By,
		LastExecutedDate: time.Now().Format(dbTSLayout),
	}

	if err := s.repo.UpdateTestStatus(status, testRunId); err != nil {
		return err
	} else {
		return nil
	}
}

func (s DefaultTestRunService) UpdateTestRun(id string, req dto.AddTestRunRequest) *errs.AppError {
	log.Info("line23")
	log.Info(req.TestCases)

	c := domain.TestRun{
		TestRun_Id:       "",
		Name:             req.Name,
		Description:      req.Description,
		Release_Id:       req.Release_Id,
		Assignee:         req.Assignee,
		TestCases:        req.TestCases,
		Status:           "Unexecuted",
		Executed_By:      req.Assignee,
		LastExecutedDate: time.Now().Format(dbTSLayout),
	}

	if err := s.repo.UpdateTestRun(id, c); err != nil {
		return err
	} else {
		return nil
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

func (s DefaultTestRunService) GetTestCaseTitlesForTestRun(id string) ([]dto.GetTestCaseTitleTestRun, *errs.AppError) {

	testRuns, err := s.repo.GetTestCaseTitlesForTestRun(id)
	if err != nil {
		return nil, err
	}
	response := make([]dto.GetTestCaseTitleTestRun, 0)
	for _, testRun := range testRuns {
		response = append(response, testRun.ToTestCaseTitleTestRunDto())
	}
	return response, err
}
func (s DefaultTestRunService) GetProjectTestRun(project_id string, id string) (*dto.GetTestRun, *errs.AppError) {

	testRun, err := s.repo.GetProjectTestRun(project_id, id)
	if err != nil {
		return nil, err
	}
	response := testRun.ToTestRunDto()
	return &response, err
}

func NewTestRunService(repository domain.TestRunRepository) DefaultTestRunService {
	return DefaultTestRunService{repository}
}
