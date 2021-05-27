package subscription

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

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/v72"
	portalsession "github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/webhook"
)

type CreateCheckoutDto struct {
	Price string `json:"priceId"`
}

type ErrResp struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type CreateCheckoutResult struct {
	SessionID string `json:"sessionId"`
}

type UpgradeSubscriptionDto struct {
	UserId    int    `json:"userId"`
	SessionId string `json:"sessionId"`
}

type ManageSubscriptionResult struct {
	URL string `json:"url"`
}

type HandleCustomerPortalDto struct {
	SessionID string `json:"sessionId"`
}

type HandleCustomerPortalResult struct {
	URL string `json:"url"`
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

	var createCheckoutDto CreateCheckoutDto
	err = json.Unmarshal(requestBody, &createCheckoutDto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(os.Getenv("SITE_URL") + "/subscription/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(os.Getenv("SITE_URL") + "/subscription/canceled"),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(createCheckoutDto.Price),
				Quantity: stripe.Int64(1),
			},
		},
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

func UpgradeSubscription(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var upgradeSubscriptionDto UpgradeSubscriptionDto
	err = json.Unmarshal(requestBody, &upgradeSubscriptionDto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.ValidUser(request, upgradeSubscriptionDto.UserId); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	err = repo.UpgradeSubscription(upgradeSubscriptionDto.UserId, upgradeSubscriptionDto.SessionId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func HandleCustomerPortal(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var handleCustomerPortalDto HandleCustomerPortalDto
	err = json.Unmarshal(requestBody, &handleCustomerPortalDto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	s, err := session.Get(handleCustomerPortalDto.SessionID, nil)
	if err != nil {
		writeJSON(writer, nil, err)
		return
	}

	returnURL := os.Getenv("SITE_URL")
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(s.Customer.ID),
		ReturnURL: stripe.String(returnURL),
	}

	ps, _ := portalsession.New(params)

	result := HandleCustomerPortalResult{
		URL: ps.URL,
	}
	writeJSON(writer, result, nil)
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
	case "customer.subscription.deleted":
		var req struct {
			Customer string `json:"customer"`
		}

		err := json.Unmarshal(event.Data.Raw, &req)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
		c, err := customer.Get(req.Customer, nil)
		if err != nil {
			http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
			return
		}

		err = repo.UnsubscribeUser(c.Email)
		if err != nil {
			http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
			return
		}
	}
}
