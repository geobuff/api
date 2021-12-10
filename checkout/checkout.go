package checkout

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/geobuff/api/repo"
	"github.com/stripe/stripe-go/client"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/webhook"
)

type ErrResp struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type CreateCheckoutResult struct {
	SessionID string `json:"sessionId"`
}

type WebhookDto struct {
	Customer string `json:"customer"`
}

func HandleCreateCheckoutSession(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var createCheckoutDto repo.CreateCheckoutDto
	err = json.Unmarshal(requestBody, &createCheckoutDto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	_, err = repo.InsertOrder(createCheckoutDto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	merch, err := repo.GetMerch()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, checkoutItem := range createCheckoutDto.Items {
		for _, merchItem := range merch {
			if checkoutItem.ID == merchItem.ID {
				amount := int64(merchItem.Price.Float64 * 100)
				image := os.Getenv("SITE_URL") + merchItem.Images[0].ImageUrl
				newItem := stripe.CheckoutSessionLineItemParams{
					Amount:   &amount,
					Name:     stripe.String(fmt.Sprintf("%s - %s", merchItem.Name, checkoutItem.SizeName)),
					Images:   []*string{stripe.String(image)},
					Currency: stripe.String("NZD"),
					Quantity: stripe.Int64(int64(checkoutItem.Quantity)),
				}
				lineItems = append(lineItems, &newItem)
			}
		}
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(os.Getenv("SITE_URL") + "/checkout/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(fmt.Sprintf("%s/checkout/canceled?email=%s", os.Getenv("SITE_URL"), createCheckoutDto.Customer.Email)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems:     lineItems,
		CustomerEmail: &createCheckoutDto.Customer.Email,
	}

	s, err := session.New(params)
	if err != nil {
		writeJSON(writer, nil, err)
		return
	}

	result := CreateCheckoutResult{
		SessionID: s.ID,
	}
	writeJSON(writer, result, nil)
}

func writeJSON(w http.ResponseWriter, v interface{}, err error) {
	var respVal interface{}
	if err != nil {
		msg := err.Error()
		var serr *stripe.Error
		if errors.As(err, &serr) {
			msg = serr.Msg
		}
		w.WriteHeader(http.StatusBadRequest)
		var e ErrResp
		e.Error.Message = msg
		respVal = e
	} else {
		respVal = v
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(respVal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}

func HandleWebhook(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	event, err := webhook.ConstructEvent(requestBody, request.Header.Get("Stripe-Signature"), os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var req struct {
			Customer string `json:"customer"`
		}

		err := json.Unmarshal(event.Data.Raw, &req)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		sc := &client.API{}
		sc.Init(os.Getenv("STRIPE_SECRET_KEY"), nil)
		c, err := sc.Customers.Get(req.Customer, nil)
		if err != nil {
			http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
			return
		}

		err = repo.UpdateStatusLatestOrder(c.Email)
		if err != nil {
			http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
			return
		}
	}
}
