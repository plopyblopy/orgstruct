package rest

import (
	"encoding/json"
	"net/http"

	"github.com/plopyblopy/orgstruct/internal/domain"
)

// RespondJSON формирует ответ в формате json
func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, err error) {
	var status int
	if sc, ok := err.(domain.HTTPStatusCoder); ok {
		status = sc.Code()
	} else {
		status = http.StatusInternalServerError
	}

	if mr, ok := err.(domain.MultiRow); ok {
		if mr.IsRow() {
			RespondRowError(w, status, err.Error())
		} else {
			RespondMultiRowError(w, status, mr.Rows())
		}
	} else {
		RespondRowError(w, status, err.Error())
	}
}

func RespondRowError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}

func RespondMultiRowError[T any](w http.ResponseWriter, status int, errs T) {
	RespondJSON(w, status, map[string]interface{}{"errors": errs})
}
