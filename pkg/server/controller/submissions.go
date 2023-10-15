package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type Submission struct {
	Code                     string   `json:"code"`
	Language                 string   `json:"lang"`
	CompilerOptions          string   `json:"compiler_options"`
	CommandLineArguments     string   `json:"command_line_arguments"`
	Stdin                    []string `json:"stdin"`
	ExpectedOutputs          []string `json:"expected_outputs"`
	CallbackUrl              string   `json:"callback_url"`
	Iterations               int32    `json:"iterations"`
	RedirectStderrToStdout   bool     `json:"redirect_stderr_to_stdout"`
	StackLimit               int32    `json:"stack_limit"`
	MemoryLimit              int32    `json:"memory_limit"`
	CpuTimeLimit             float32  `json:"cpu_time_limit"`
	CpuExtraTime             float32  `json:"cpu_extra_time"`
	WallTimeLimit            float32  `json:"wall_time_limit"`
	MaxProcessesAndOrThreads int32    `json:"max_processes_and_or_threads"`
}

func WriteError(w http.ResponseWriter, code int, reason string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(
		fmt.Sprintf(`{ "reason": %q, "statusCode": %d }`, reason, code)),
	)
}

func HandleSubmissionsPost(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		WriteError(w, http.StatusBadRequest, "invalid content-type")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var submission Submission = Submission{
		Language: "cpp",
	}

	err := decoder.Decode(&submission)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(submission.Code) == 0 {
		WriteError(w, http.StatusBadRequest, "source code is missing")
	}

	if len(submission.Language) == 0 {
		WriteError(w, http.StatusBadRequest, "language field is missing")
	}
}

func HandleSubmissionsGetOne(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	w.Write([]byte(token))
}

func HandleSubmissionsGetMany(w http.ResponseWriter, r *http.Request) {
	// page := r.URL.Query().Get("page")
	// perPage := r.URL.Query().Get("perPage")
	// fields := r.URL.Query().Get("fields")
}
