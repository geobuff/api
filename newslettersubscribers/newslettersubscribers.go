package newslettersubscribers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func CreateNewsletterSubscriber(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var newSubscriber repo.CreateNewsletterSubscriberDto
	err = json.Unmarshal(requestBody, &newSubscriber)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	exists, err := repo.NewsletterSubscriberExists(newSubscriber.Email)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(writer, fmt.Sprintf("The email address %s is already subscribed to our newsletter.", newSubscriber.Email), http.StatusBadRequest)
		return
	}

	if err := repo.InsertNewsletterSubscriber(newSubscriber); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func Unsubscribe(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	err = repo.DeleteNewsletterSubscriber(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
