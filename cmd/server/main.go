package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/srijanmukherjee/codesensei/pkg/config"
	"github.com/srijanmukherjee/codesensei/pkg/server/controller"
)

func main() {
	router := chi.NewRouter()

	router.Route("/submissions", func(r chi.Router) {
		r.Post("/", controller.HandleSubmissionsPost)
		r.Get("/", controller.HandleSubmissionsGetMany)
		r.Get("/{token}", controller.HandleSubmissionsGetOne)
	})

	addr := fmt.Sprintf("0.0.0.0:%v", config.Config.ServerPort)
	log.Printf("🚀 server started on http//%v\n", addr)
	log.Panic(http.ListenAndServe(addr, router))
}
