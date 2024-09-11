package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type PingError struct {
	URL        string
	StatusCode int
}

func (e *PingError) Error() string {
	return fmt.Sprintf("failed to ping service: %d", e.StatusCode)
}

type HTTPChecker struct {
	client *http.Client
	URL    string
}

func NewHTTPChecker(url string, c *http.Client) *HTTPChecker {
	return &HTTPChecker{client: c, URL: url}
}

func (c *HTTPChecker) Ping(ctx context.Context) error {
	fmt.Printf("pinging : %s\n", c.URL)
	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, c.URL, nil)
	if reqErr != nil {
		return reqErr
	}
	res, resErr := c.client.Do(req)
	if resErr != nil {
		return resErr
	}
	if res.StatusCode > 399 {
		return &PingError{URL: c.URL, StatusCode: res.StatusCode}
	}
	fmt.Printf("success : %s\n", c.URL)
	return nil
}

type GRPCChecker struct {
	endpoint string
	client   grpc_health_v1.HealthClient
}

func NewGRPCChecker(endpoint string) *GRPCChecker {
	conn, connErr := grpc.NewClient(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if connErr != nil {
		panic(connErr)
	}
	return &GRPCChecker{
		client:   grpc_health_v1.NewHealthClient(conn),
		endpoint: endpoint,
	}
}

func (c *GRPCChecker) Ping(ctx context.Context) error {
	// Implement gRPC ping logic here
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	resp, err := c.client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return err
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return &PingError{
			URL:        c.endpoint,
			StatusCode: int(resp.Status),
		}
	}
	return nil
}
