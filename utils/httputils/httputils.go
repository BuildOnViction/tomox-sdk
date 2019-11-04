package httputils

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, code int, message string) {
	Write(w, code, map[string]string{"error": message})
}
func WriteMessage(w http.ResponseWriter, code int, message string) {
	Write(w, code, map[string]string{"message": message})
}
func WriteJSON(w http.ResponseWriter, code int, payload interface{}) {
	Write(w, code, map[string]interface{}{"data": payload})
}

func Write(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
