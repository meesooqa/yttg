package web

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/meesooqa/yttg/app/job"
	"github.com/meesooqa/yttg/app/web/mocks"
)

func TestRouterStatic(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("static content"), 0644)
	assert.NoError(t, err)

	server := &Server{
		TemplStaticLocation: tmpDir,
		templates:           template.Must(template.New("").Parse("")),
	}
	router := server.router()

	req, err := http.NewRequest("GET", "/static/test.txt", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "static content", rr.Body.String())
}

func TestRouterRoot(t *testing.T) {
	mockQueue := new(mocks.MockJobQueue)
	mockQueue.On("GetJobsStatuses").Return(map[string]job.JobStatus{
		"001": job.StatusQueued,
	})

	server := &Server{
		JobQueue:  mockQueue,
		templates: template.Must(template.New("index").Parse("")),
	}
	router := server.router()

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockQueue.AssertCalled(t, "GetJobsStatuses")
}

func TestRouterSendPost(t *testing.T) {
	mockQueue := new(mocks.MockJobQueue)
	mockQueue.On("AddJob", mock.AnythingOfType("job.SendVideoJob")).Return()

	server := &Server{
		JobQueue:  mockQueue,
		templates: template.Must(template.New("test").Parse("")),
	}
	router := server.router()

	form := url.Values{}
	form.Add("url", "http://example.com")
	req, err := http.NewRequest("POST", "/send", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusSeeOther, rr.Code)
	mockQueue.AssertCalled(t, "AddJob", mock.AnythingOfType("job.SendVideoJob"))
}
