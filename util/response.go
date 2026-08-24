package util

import (
	"encoding/json"
	"net/http"
)

func RespondWithErrorFromCode(w http.ResponseWriter, code int) {

	http.Error(w, http.StatusText(code), code)
}

func RespondWithErrorCustom(w http.ResponseWriter, custom_err string, code int) {

	http.Error(w, custom_err, code)

}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
