package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/albionmaksbad-jpg/distributed-job-queue/internal/job"
)

type testStore struct {
	job job.Job
	err error
	in  job.CreateInput
}

func (s *testStore) Create(_ context.Context, in job.CreateInput) (job.Job, error) {
	s.in = in
	return s.job, s.err
}

func TestCreateJob(t *testing.T) {
	s := &testStore{job: job.Job{ID: "job-1", Type: "report", Status: job.StatusPending}}
	req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(`{"type":"report","payload":{"format":"csv"}}`))
	res := httptest.NewRecorder()
	NewHandler(s).ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusCreated)
	}
	var got job.Job
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.ID != s.job.ID || s.in.Type != "report" {
		t.Fatalf("response = %+v, input = %+v", got, s.in)
	}
}

func TestCreateJobRejectsInvalidBody(t *testing.T) {
	cases := []string{
		`{"type":"report"}`,
		`{"type":"report","payload":{},"extra":true}`,
		`{"type":"report","payload":{}} {}`,
	}
	for _, body := range cases {
		t.Run(body, func(t *testing.T) {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(body))
			NewHandler(&testStore{}).ServeHTTP(res, req)
			if res.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestCreateJobStoreError(t *testing.T) {
	s := &testStore{err: errors.New("database unavailable")}
	req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(`{"type":"report","payload":{}}`))
	res := httptest.NewRecorder()
	NewHandler(s).ServeHTTP(res, req)
	if res.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusInternalServerError)
	}
}
