package web

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"

	"task-queue-001/app/job"
	"task-queue-001/app/media"
	"task-queue-001/app/send"
)

func (s *Server) getIndexPageCtrl(w http.ResponseWriter, r *http.Request) {
	statuses := s.JobQueue.GetJobsStatuses()
	s.templates.Execute(w, &statuses)
}

func (s *Server) getStatusPageCtrl(w http.ResponseWriter, r *http.Request) {
	statuses := s.JobQueue.GetJobsStatuses()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		http.Error(w, "JSON Encoding error", http.StatusInternalServerError)
	}
}

func (s *Server) send(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")
		if url == "" {
			http.Error(w, "URL can't be empty", http.StatusBadRequest)
			return
		}
		s.JobQueue.AddJob(job.SendVideoJob{
			BaseJob:         job.BaseJob{ID: uuid.New().String()},
			URL:             url,
			MediaService:    media.NewMediaService(),
			TelegramFactory: &send.EnvClientFactory{},
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
}
