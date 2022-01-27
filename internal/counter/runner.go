package counter

import (
	"context"
	"net/http"

	api "challenge/pkg/api/http"
	"challenge/pkg/counter"
	"challenge/pkg/counter/memory"
)

type Service struct {
	// API server internal data
	cfg        Config
	mux        *http.ServeMux
	httpServer *http.Server
	counterSvc counter.Counter
}

func NewService(cfg Config) (svc *Service) {
	svc = new(Service)
	// process config
	svc.cfg = cfg
	svc.mux = http.NewServeMux()
	svc.counterSvc = memory.NewCounter()
	// setup counter API on http request muxer
	api.NewCounterAPI(svc.counterSvc).Setup(svc.mux)
	handler := http.Handler(svc.mux)
	// activate cors if config enables them
	if svc.cfg.Cors {
		handler = api.SetupCorsMiddleware(handler)
	}
	// http server instantiation
	svc.httpServer = &http.Server{Addr: cfg.Listen, Handler: handler}
	return
}

func (svc *Service) Run() error {
	return svc.httpServer.ListenAndServe()
}

func (svc *Service) Shutdown(ctx context.Context) error {
	return svc.httpServer.Shutdown(ctx)
}
