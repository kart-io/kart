package kart_grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kart-io/kart/internal/host"
)

type testKey struct{}

func TestServer(t *testing.T) {
	config := &GrpcConfig{
		Port: "8081",
		Addr: "0.0.0.0",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, testKey{}, "test")
	srv := NewGrpcServer(WithConfig(config))

	if e, err := srv.Endpoint(); err != nil || e == nil {
		t.Fatal(e, err)
	}

	go func() {
		// start server
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second * 10)
	testClient(t, srv)
	if err := srv.Stop(ctx); err != nil {
		t.Fatal(err)
	}
}

func testClient(t *testing.T, srv *GrpcServer) {
	port, ok := host.Port(srv.listener)
	if !ok {
		t.Fatalf("extract port error: %v", srv.listener)
	}
	endpoint := fmt.Sprintf("127.0.0.1:%d", port)
	// new a gRPC client
	conn, err := DialInsecure(context.Background(), WithEndpoint(endpoint))
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}
