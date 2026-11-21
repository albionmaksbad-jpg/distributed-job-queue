package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/albionmaksbad-jpg/distributed-job-queue/internal/job"
)

const maxBodySize = 1 << 20

type JobStore interface {
	Create(context.Context, job.CreateInput) (job.Job, error)
}

func NewHandler(store JobStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", createJob(store))
	return mux
}

func createJob(store JobStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var in job.CreateInput
		if err := dec.Decode(&in); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := in.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		j, err := store.Create(r.Context(), in)
		if err != nil {
			http.Error(w, "could not create job", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(j)
	}
}
