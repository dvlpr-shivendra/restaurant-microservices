package discovery

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
	addressList, err := registry.Discover(ctx, serviceName)

	if err != nil {
		return nil, err
	}

	return grpc.NewClient(addressList[rand.Intn(len(addressList))], grpc.WithTransportCredentials(insecure.NewCredentials()))
}
