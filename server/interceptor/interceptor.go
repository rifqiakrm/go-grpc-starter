package interceptor

import (
	"context"

	"cloud.google.com/go/errorreporting"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorReporting reports error to Google Cloud Error Reporting.
// Only error with codes.Unknown and codes.Internal that are sent.
func ErrorReporting(client *errorreporting.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			code := status.Code(err)
			if code == codes.Unknown || code == codes.Internal {
				client.Report(errorreporting.Entry{
					Error: err,
				})
				return resp, status.Error(code, "")
			}
		}
		return resp, err
	}
}
