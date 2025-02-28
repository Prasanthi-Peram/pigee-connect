package main

import(
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
type application struct{
	config config
}

type config struct{
	addr string
}

func (app *application) mount() *chi.Mux{

	r := chi.NewRouter()
	r.Use(middleware.Logger)

		//Group of Routers
	r.Route("/v1",func(r chi.Router){
		//Sub-Router
		r.Get("/health", app.healthCheckHandler)
	})
	return r
}

func (app *application) run(mux http.Handler) error{
	srv:=&http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second*30,
		ReadTimeout: time.Second*10,
		IdleTimeout: time.Minute,
	}
	log.Printf("server has started at %s",app.config.addr)

	return srv.ListenAndServe()
}