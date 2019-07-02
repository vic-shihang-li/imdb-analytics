package job

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ExtractionJob describes a TVSeries data extraction job.
// There are 3 possible statuses: Ready | Processing | Completed
type ExtractionJob struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Handler encapsulates input and output channels for TVJob
type Handler struct {
	in  chan interface{}
	out chan<- *ExtractionJob
}

// Jobs container
var jobs []ExtractionJob

// Routes return a router with routes associated with TVSeriesExtractionJobs
func Routes(in chan interface{}, out chan<- *ExtractionJob) *mux.Router {
	r := mux.NewRouter()
	h := &Handler{in: in, out: out}
	r.HandleFunc("/", h.homeHandler)
	r.HandleFunc("/jobs", h.getJobs).Methods("GET")
	r.HandleFunc("/jobs/{id}", h.getJob).Methods("GET")
	r.HandleFunc("/jobs", h.postJob).Methods("POST")

	return r
}

// homeHandler responds a helper message that directs client to use
// the actual API at `/jobs`.
func (h *Handler) homeHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&map[string]interface{}{
		"Message": "This is the home route of the extractor service API. " +
			"To make RESTful calls visit `/jobs.`",
	})
}

// getJob handles querying job with a specific ID.
func (h *Handler) getJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return
	}
	h.in <- id
	var out interface{} = <-h.in
	job, typecheckOk := out.(ExtractionJob)
	if !typecheckOk {
		json.NewEncoder(w).Encode(&map[string]interface{}{
			"Message": "No job exists yet.",
		})
	} else {
		json.NewEncoder(w).Encode(job)
	}
}

// getJobs reponds all jobs in the system.
func (h *Handler) getJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if len(jobs) == 0 {
		json.NewEncoder(w).Encode(&map[string]interface{}{
			"Message": "No job exists yet.",
		})
	} else {
		json.NewEncoder(w).Encode(jobs)
	}
}

// postJob creates a new job, if a job with the same Name does not already exist.
func (h *Handler) postJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if name, ok := r.URL.Query()["name"]; !ok {

	} else {
		for _, item := range jobs {
			if item.Name == name[0] {
				json.NewEncoder(w).Encode(item)
				return
			}
		}
		j := ExtractionJob{
			ID:     rand.Intn(1000000000),
			Name:   name[0],
			Status: "Ready",
		}
		jobs = append(jobs, j)
		h.out <- &j
		json.NewEncoder(w).Encode(&j)
	}
}