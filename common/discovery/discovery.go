package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	Register(ctx context.Context, instanceId, serverName, hostPort string) error
	DeRegister(ctx context.Context, instanceId, serverName string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceId, serviceName string) error
}

func GenerateInstanceId(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
