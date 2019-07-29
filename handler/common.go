package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "identity")
	w.WriteHeader(status)
	w.Write(response)
}

func RespondJSON2(w http.ResponseWriter, status int) {
	//fmt.Println("SAved!!!")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte("{}"))
}

func respondJSONforInt(w http.ResponseWriter, status int, payload interface{}) {
	buffer := bytes.NewBufferString("{")
	jsonValue, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	buffer.WriteString(fmt.Sprintf("\"avg\": %s.0", string(jsonValue)))
	buffer.WriteString("}")
	w.WriteHeader(status)
	w.Write(buffer.Bytes())
}

func RespondError(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}
