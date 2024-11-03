package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	pb "restaurant-backend/common/api"
)

type TelemetryMiddleware struct {
	next PaymentsService
}

func NewTelemetryMiddleware(next PaymentsService) PaymentsService {
	return &TelemetryMiddleware{next}
}

func (s *TelemetryMiddleware) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CreatePayment: %v", o))

	return s.next.CreatePayment(ctx, o)
}
