package stripe

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"errors"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

type StripeService struct {
	SuccessURL *string
	CancelURL  *string
}

func New(secretKey, successURL, cancelURL string) *StripeService {
	stripe.Key = secretKey
	return &StripeService{
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
	}
}

func (s *StripeService) Checkout(order domain.Order, customer domain.Customer) (string, error) {
	if len(order.Cart.Items) == 0 {
		return "", errors.New("order must contain at least one item")
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          toLineItems(order.Cart.Items),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         s.SuccessURL,
		CancelURL:          s.CancelURL,
		AutomaticTax:       &stripe.CheckoutSessionAutomaticTaxParams{Enabled: stripe.Bool(true)},
	}

	if customer.StripeID != "" {
		params.Customer = stripe.String(customer.StripeID)
	}

	session, err := session.New(params)
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
