package response

import (
	"encoding/json"
	"net/http"
	"restapi-native-go/internal/utils/errors"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Error(w http.ResponseWriter, err error) {
	var statusCode int
	var errorResp interface{}

	if appErr, ok := err.(errors.AppError); ok {
		statusCode = appErr.StatusCode
		errorResp = appErr
	} else {
		statusCode = http.StatusInternalServerError
		errorResp = map[string]string{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Terjadi kesalahan pada server",
		}
	}

	JSON(w, statusCode, errorResp)
}
