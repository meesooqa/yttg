package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/meesooqa/yttg/app/job"
	"github.com/meesooqa/yttg/app/web/mocks"
)

func TestSendHandlerPost(t *testing.T) {
	mockQueue := new(mocks.MockJobQueue)
	server := &Server{
		JobQueue:  mockQueue,
		templates: template.Must(template.New("test").Parse("")),
	}

	form := url.Values{}
	form.Add("url", "http://example.com")
	req, err := http.NewRequest("POST", "/send", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.send)

	mockQueue.On("AddJob", mock.AnythingOfType("job.SendVideoJob")).Return()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusSeeOther, rr.Code)
	assert.Equal(t, "/", rr.Header().Get("Location"))
	mockQueue.AssertCalled(t, "AddJob", mock.AnythingOfType("job.SendVideoJob"))
	addedJob := mockQueue.Calls[0].Arguments[0].(job.SendVideoJob)
	assert.Equal(t, "http://example.com", addedJob.URL)
	assert.NotEmpty(t, addedJob.ID)
}

func TestSendHandlerGet(t *testing.T) {
	server := &Server{
		templates: template.Must(template.New("test").Parse("")),
	}

	req, err := http.NewRequest("GET", "/send", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.send)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Contains(t, rr.Body.String(), "Method is not allowed")
}

func TestGetIndexPageCtrl(t *testing.T) {
	mockQueue := new(mocks.MockJobQueue)
	mockQueue.On("GetJobsStatuses").Return(map[string]job.JobStatus{
		"001": job.StatusQueued,
	})

	tmpl := template.Must(template.New("index").Parse("Jobs: {{len .}}"))
	server := &Server{
		JobQueue:  mockQueue,
		templates: tmpl,
	}

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getIndexPageCtrl)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Jobs: 1", rr.Body.String())
	mockQueue.AssertCalled(t, "GetJobsStatuses")
}

func TestSendHandlerEmptyURL(t *testing.T) {
	server := &Server{
		templates: template.Must(template.New("test").Parse("")),
	}

	form := url.Values{}
	form.Add("url", "")
	req, err := http.NewRequest("POST", "/send", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.send)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "URL can't be empty")
}

func TestGetStatusPageCtrl(t *testing.T) {
	// Создаем тестовую очередь заданий
	mockQueue := new(mocks.MockJobQueue)
	mockQueue.On("GetJobsStatuses").Return(map[string]job.JobStatus{})

	// Создаем объект Server с тестовой очередью
	s := &Server{JobQueue: mockQueue}

	// Создаем фиктивный HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	// Создаем фиктивный ResponseWriter
	rr := httptest.NewRecorder()

	// Вызываем тестируемый метод
	s.getStatusPageCtrl(rr, req)

	// Проверяем код ответа
	resp := rr.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	// Проверяем заголовок Content-Type
	if ct := resp.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Ожидался Content-Type 'application/json; charset=utf-8', получен '%s'", ct)
	}

	// Декодируем JSON-ответ
	var statuses map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&statuses); err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	// Здесь можно выполнить дополнительные проверки содержимого,
	// например, если очередь пуста, ожидается пустая карта.
	if len(statuses) != 0 {
		t.Errorf("Ожидалась пустая карта статусов, получено: %v", statuses)
	}
}
