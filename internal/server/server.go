package srv

import (
	"context"
	"log"
	"net/http"
	cfg "org-structure-api/internal/config"
	depSvc "org-structure-api/internal/core/department/service"
	empSvc "org-structure-api/internal/core/employee/service"
	"org-structure-api/internal/server/api"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg cfg.ServerConfig, depSvc depSvc.Interface, empSvc empSvc.Interface) *Server {
	r := mux.NewRouter()

	depHand, empHand := api.SetupAllHandlers(depSvc, empSvc)
	api.SetupAllRoutes(r, depHand, empHand)

	s := &Server{
		httpServer: &http.Server{
			Addr:         cfg.Addr,
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	return s
}

func (s *Server) Start() {
	go func() {
		log.Printf("Server running on %s\n", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	s.Shutdown()
}

func (s *Server) Shutdown() {
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
