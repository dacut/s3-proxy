package middlewares

import (
	"net/http"

	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/bucket"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/config"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/s3client"
	"golang.org/x/net/context"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

var bucketRequestContextKey = &contextKey{name: "bucket-request-context"}

func BucketRequestContext(
	tgt *config.TargetConfig,
	path string,
	s3clientManager s3client.Manager,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Generate new bucket client
			brctx := bucket.NewClient(tgt, path, s3clientManager)
			// Add bucket structure to request context by creating a new context
			ctx := context.WithValue(req.Context(), bucketRequestContextKey, brctx)
			// Create new request with new context
			req = req.WithContext(ctx)
			// Next
			next.ServeHTTP(rw, req)
		})
	}
}

func GetBucketRequestContext(req *http.Request) bucket.Client {
	res, _ := req.Context().Value(bucketRequestContextKey).(bucket.Client)

	return res
}
