package main

import (
	"fmt"
	"kozo/utils"
	"kozo/views/rest"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
)

func MainRestServer(env string) {

	// Load Env
	err := godotenv.Load(env)

	if err != nil {
		panic(err)
	}

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

	serverIP := os.Getenv("serverIP");
	serverPort := os.Getenv("serverPort");
	fmt.Println("REST server is up and running on Port:", serverPort)
	http.ListenAndServe(fmt.Sprintf("%s:%s", serverIP, serverPort), r)
}