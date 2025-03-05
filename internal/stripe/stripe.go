package stripe

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/models"
	"encoding/json"
	"errors"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/webhook"
)

type StripeService struct {
	successURL    string
	cancelURL     string
	webhookSecret string
	models        *models.Model
}

func New(secretKey, successURL, cancelURL, webhookSecret string, model *models.Model) *StripeService {
	stripe.Key = secretKey
	return &StripeService{
		successURL:    successURL,
		cancelURL:     cancelURL,
		webhookSecret: webhookSecret,
		models:        model,
	}
}

func (s *StripeService) ConfirmOrder(payload []byte, header string) error {
	event, err := webhook.ConstructEvent(payload, header, s.webhookSecret)
	if err != nil {
		return err
	}

	var session stripe.CheckoutSession

	err = json.Unmarshal(event.Data.Raw, &session)
	if err != nil {
		return err
	}

	if session.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
		return errors.New("payment has failed")
	}

	order, err := s.models.Orders.GetByStripeID(session.ID)
	if err != nil {
		return err
	}

	err = s.models.Orders.Update(order.ID, domain.StatusCompleted)
	if err != nil {
		return err
	}

	return nil
}

func (s *StripeService) Checkout(order domain.Order, customer domain.Customer) (string, error) {
	if len(order.Cart.Items) == 0 {
		return "", errors.New("order must contain at least one item")
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          toLineItems(order.Cart.Items),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(s.successURL),
		CancelURL:          stripe.String(s.cancelURL),
		AutomaticTax:       &stripe.CheckoutSessionAutomaticTaxParams{Enabled: stripe.Bool(true)},
	}

	if customer.StripeID != "" {
		params.Customer = stripe.String(customer.StripeID)
	}

	session, err := session.New(params)
	if err != nil {
		return "", err
	}

	err = s.models.Orders.AddStripeID(order.ID, session.ID)
	if err != nil {
		return "", err
	}

	return session.URL, nil
}

func toLineItem(item domain.CartItem) *stripe.CheckoutSessionLineItemParams {
	return &stripe.CheckoutSessionLineItemParams{
		PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
			Currency: stripe.String("egp"),
			ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
				Name:        stripe.String(item.Name),
				Description: stripe.String(item.Description),
				Metadata:    item.Properties,
				Images:      stripe.StringSlice(item.Images),
			},
			UnitAmount: stripe.Int64(item.ItemPrice.Amount()),
		},
		Quantity: stripe.Int64(int64(item.Amount)),
	}
}

func toLineItems(items []domain.CartItem) []*stripe.CheckoutSessionLineItemParams {
	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, item := range items {
		lineItem := toLineItem(item)

		lineItems = append(lineItems, lineItem)
	}

	return lineItems
}
