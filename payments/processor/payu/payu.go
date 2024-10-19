package payu

import (
	pb "restaurant-backend/common/api"
)

type payu struct{}

// NewProcessor initializes a new PayU payment processor.
func NewProcessor() *payu {
	return &payu{}
}

// CreatePaymentLink generates a payment link for the given order using PayU's API.
func (p *payu) CreatePaymentLink(order *pb.Order) (string, error) {
	// TODO: Implement PayU payment link generation logic
	return "https://apitest.payu.in/public/#/7243ac088aa7849da7753d58bf4bdc26b8929f989dbd50d7bb631208618e61ba", nil
}
