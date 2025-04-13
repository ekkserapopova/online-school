package responser

import (
	"encoding/json"
	"net/http"
)

type MessageResponse struct {
	Msg string `json:"msg"`
}

func SendOk(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp, err := json.Marshal(v)
	if err != nil {
		return
	}
	_, _ = w.Write(resp)
}

func SendErr(w http.ResponseWriter, code int, msg string) {
	resp, err := json.Marshal(MessageResponse{msg})
	if err != nil {
		return
	}
	w.WriteHeader(code)
	_, _ = w.Write(resp)
}
