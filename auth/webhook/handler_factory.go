package webhook

import (
	"net/http"

	"go.containerssh.io/containerssh/internal/auth"
	"go.containerssh.io/containerssh/log"
)

// NewHandler creates an HTTP handler that forwards calls to the provided h config request handler.
func NewHandler(h AuthRequestHandler, logger log.Logger) http.Handler {
	return auth.NewHandler(h, logger)
}
