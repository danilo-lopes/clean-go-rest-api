// Clean Architecture - Frameworks & Drivers Layer
// HTTP Server setup
package server

import (
	"clean-go-rest-api/internal/crosscutting/logger"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer(router *mux.Router, port int, logger logger.ILogger) {
	logger.Info(fmt.Sprintf("Server running on port %d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		log.Fatal(err)
	}
}
