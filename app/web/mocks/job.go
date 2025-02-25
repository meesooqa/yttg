package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/meesooqa/yttg/app/job"
)

type MockJobQueue struct {
	mock.Mock
}

func (m *MockJobQueue) AddJob(j job.Job) {
	m.Called(j)
}

func (m *MockJobQueue) GetJobsStatuses() map[string]job.JobStatus {
	args := m.Called()
	return args.Get(0).(map[string]job.JobStatus)
}
