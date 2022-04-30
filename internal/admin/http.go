package admin

import (
	"context"
	"log"
	"net"
	"net/http"
)

type server struct {
	config          *Config
	addr            string
	server          *http.Server
	admin           *adminHandler
	idleConnsClosed chan struct{}
}

func NewServer(cfg *Config) *server {
	return &server{
		addr:            net.JoinHostPort(cfg.Host, cfg.Port),
		config:          cfg,
		admin:           newAdminHandler(cfg, defaultRoutes(cfg)),
		idleConnsClosed: nil,
	}
}

func (svr *server) Shutdown() {
	err := svr.server.Shutdown(context.Background())
	if err != nil {
		log.Printf("server shutdown: %v", err)
	}
	close(svr.idleConnsClosed)
}

func (svr *server) Serve() (err error) {
	svr.setupServer()
	svr.setupHandler()

	svr.idleConnsClosed = make(chan struct{})

	log.Printf("admin server is running at \"http://%s\"\n", svr.addr)

	err = svr.server.ListenAndServe()
	if err != http.ErrServerClosed {
		err = nil
	}
	if err != nil {
		return err
	}

	<-svr.idleConnsClosed

	log.Printf("server closed")

	return
}

func (svr *server) setupServer() {
	svr.server = &http.Server{
		Addr: svr.addr,
	}
}

func (svr *server) setupHandler() {
	svr.admin.addRoute(http.MethodPost, `^/shutdown$`, svr.handleShutdown)
	handler := requestLogger(svr.admin)
	if svr.config.IsBasicAuth {
		handler = basicAuth(handler,
			svr.config.BasicAuthUsername,
			svr.config.BasicAuthPassword)
	}
	svr.server.Handler = handler
}

func (svr *server) handleShutdown(w http.ResponseWriter, r *http.Request) {
	log.Println("received shutdown request")
	w.Write([]byte("shutting down the dagman server...\n"))
	go func() {
		svr.Shutdown()
	}()
}
