package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// Rest is responsible to act as HTTP/1.1 REST server.
// It composes grpc-gateway runtime.ServeMux.
type Rest struct {
	*runtime.ServeMux
	port string
}

// ErrorData represents error code and message for the response.
type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error represents error response.
type Error struct {
	// Errors represents list of errors that are visible in response.
	Error ErrorData `json:"error"`
	// Meta represents auxiliary data that is visible in response.
	Meta interface{} `json:"meta"`
}

// NewRest creates an instance of Rest.
func NewRest(port string) *Rest {
	return &Rest{
		ServeMux: runtime.NewServeMux(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			}),
			runtime.WithErrorHandler(customErrorHandler),
		),
		port: port,
	}
}

// NewProductionRest creates an instance of Rest with default production options attached.
// The only difference between NewRest and NewProductionRest is the later enable Prometheus metrics by default.
func NewProductionRest(port string) *Rest {
	srv := &Rest{
		ServeMux: runtime.NewServeMux(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			}),
			runtime.WithErrorHandler(customErrorHandler),
		),
		port: port,
	}
	_ = srv.EnablePrometheus() // error is impossible, hence ignored.
	_ = srv.EnableHealth()     // error is impossible, hence ignored.
	return srv
}

// EnablePrometheus enables prometheus endpoint.
// It can be accessed via /metrics.
func (r *Rest) EnablePrometheus() error {
	return r.ServeMux.HandlePath(http.MethodGet, "/metrics", prometheusHandler())
}

// EnableHealth enables health endpoint.
// It can be accessed via /healthz.
func (r *Rest) EnableHealth() error {
	return r.ServeMux.HandlePath(http.MethodGet, "/healthz", healthHandler())
}

// Run runs HTTP/1.1 runtime.ServeMux.
// It runs inside a goroutine.
func (r *Rest) Run() error {
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", r.port), allowCORS(r.ServeMux)); err != nil {
			panic(err)
		}
	}()
	return nil
}

func prometheusHandler() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	}
}

func healthHandler() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.WriteHeader(http.StatusOK)
	}
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				headers := []string{"Content-Type", "Accept", "Authorization"}
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
				methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// customErrorHandler customize error handler instead of using the default one
func customErrorHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	mrs runtime.Marshaler,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	// return default error message when marshal failed
	const fallback = `{"error": {"code":13,"message":"failed to marshal error message"}, "meta":null}`

	s := status.Convert(err)

	w.Header().Set("Content-type", mrs.ContentType("application/json"))
	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	jsonErr := json.NewEncoder(w).Encode(Error{
		Error: ErrorData{
			Code:    int(s.Code()),
			Message: s.Message(),
		},
		Meta: nil,
	})

	if jsonErr != nil {
		_, _ = w.Write([]byte(fallback))
	}
}
