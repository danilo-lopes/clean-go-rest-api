// Clean Architecture - Frameworks & Drivers Layer
// HTTP Server setup
package server

import (
	"clean-go-rest-api/internal/crosscutting/logger"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer(router *mux.Router, port int, logger logger.ILogger) *http.Server {
	addr := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("HTTP server ListenAndServe: %v", err))
		}
	}()
	return httpServer
}
