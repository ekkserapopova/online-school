package reader

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func ReadRequestData(r *http.Request, request interface{}) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if err = json.Unmarshal(data, &request); err != nil {
		return err
	}
	return nil
}

func ReadVarsUUID(r *http.Request, key string) (uuid.UUID, error) {
	vars := mux.Vars(r)
	uuidStr := vars[key]
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
