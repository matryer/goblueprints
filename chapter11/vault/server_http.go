package vault

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

// NewHTTPServer makes a new Vault HTTP service.
func NewHTTPServer(ctx context.Context, vaultService Service) http.Handler {
	m := http.NewServeMux()
	m.Handle("/hash", httptransport.NewServer(
		ctx,
		makeHashEndpoint(vaultService),
		decodeHashRequest,
		encodeResponse,
	))
	m.Handle("/validate", httptransport.NewServer(
		ctx,
		makeValidateEndpoint(vaultService),
		decodeValidateRequest,
		encodeResponse,
	))
	return m
}
