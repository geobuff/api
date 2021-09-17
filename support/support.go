package support

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/geobuff/api/email"
)

func SendSupportRequest(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var supportRequest email.SupportRequest
	err = json.Unmarshal(requestBody, &supportRequest)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	_, err = email.SendSupportRequest(supportRequest)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
