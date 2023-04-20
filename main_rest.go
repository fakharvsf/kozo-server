package main

import (
	"fmt"
	"kozo/utils"
	"kozo/views/rest"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func MainRestServer() {
	// DB Connection
	dbError := utils.DBConnect()

	if dbError != nil {
		panic(dbError)
	}

	// Run migrations
	utils.DBMigrate(false)

	// Router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res := utils.ARSuccess("Server is up and running.")
		render.JSON(w, r, res)
	})

	r.Route("/api/auth", rest.Auth)
	r.Route("/api/tasks", rest.Tasks)
	r.Route("/api/users", rest.Users)

	RestServer := "192.168.1.4";
	// RestServer := "127.0.0.1";
	RestPORT := 3333;
	fmt.Println("REST server is up and running on Port:", RestPORT)
	http.ListenAndServe(fmt.Sprintf("%s:%d", RestServer, RestPORT), r)
}