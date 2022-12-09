package server

import (
	"fmt"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// DialOption allows optional config for dialer
type DialOption func(name string) (grpc.DialOption, error)

// WithTracer traces rpc calls
func WithTracer(t opentracing.Tracer) DialOption {
	return func(name string) (grpc.DialOption, error) {
		return grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(t)), nil
	}
}

// Dial returns a load balanced grpc client conn with tracing interceptor
func Dial(name string, opts ...DialOption) (*grpc.ClientConn, error) {
	dialopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	for _, fn := range opts {
		opt, err := fn(name)
		if err != nil {
			return nil, fmt.Errorf("config error: %v", err)
		}
		dialopts = append(dialopts, opt)
	}

	conn, err := grpc.Dial(name, dialopts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %v", name, err)
	}

	return conn, nil
}

// DialWithSSL returns a load balanced grpc client conn with tracing interceptor
func DialWithSSL(name string, certFile string, opts ...DialOption) (*grpc.ClientConn, error) {
	cred, errSSL := credentials.NewClientTLSFromFile(certFile, "")
	if errSSL != nil {
		log.Fatalf("error while reading cert file : %v", errSSL)
	}

	dialopts := []grpc.DialOption{
		grpc.WithTransportCredentials(cred),
	}

	for _, fn := range opts {
		opt, err := fn(name)
		if err != nil {
			return nil, fmt.Errorf("config error: %v", err)
		}
		dialopts = append(dialopts, opt)
	}

	conn, err := grpc.Dial(name, dialopts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %v", name, err)
	}

	return conn, nil
}

// InitGRPCConn returns gRPC client connection for connecting to another service
func InitGRPCConn(addr string, ssl bool, cert string) *grpc.ClientConn {
	if ssl {
		conn, err := DialWithSSL(addr, cert)
		if err != nil {
			panic(fmt.Sprintf("ERROR: dial error: %v", err))
		}
		return conn
	}

	conn, err := Dial(addr)
	if err != nil {
		panic(fmt.Sprintf("ERROR: dial error: %v", err))
	}
	return conn
}
