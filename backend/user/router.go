package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func Init(r *mux.Router) {
	r.HandleFunc("/user/register", registerUser)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Invalid request body"))
	}

	var body map[string]string
	err = json.Unmarshal(reqBody, &body)
	if err != nil {
		w.Write([]byte("Invalid request body"))
	}

	username, ok := body["username"]
	if !ok {
		w.Write([]byte("Invalid request body. Missing username parameter."))
	}

	username = strings.Trim(username, " \n")

	if len(username) < 3 {
		w.Write([]byte("Invalid username parameter."))
	}

	user, err := saveUser(username)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	response := make(map[string]string, 0)
	response["token"] = user.Id
	data, _ := json.Marshal(response)
	w.Write(data)
}
