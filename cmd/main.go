package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"rss-agg/internal/database"
	"rss-agg/internal/handler"
	"rss-agg/internal/middleware"
	"rss-agg/internal/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Running...")

	// load env data
	godotenv.Load("../.env")
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT is empty")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("dbURL is empty")
	}
	// create connection ot DB
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	dbConn := database.New(conn)

	go utils.StartScraping(dbConn, 10, time.Minute)

	mainRouter := chi.NewRouter()
	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1ApiRouter := chi.NewRouter()
	v1ApiRouter.Use(middleware.DBConn(dbConn))

	// public group
	v1ApiRouter.Group(func(public chi.Router) {
		public.Get("/healthz", handler.HealthCheck)
		public.Get("/err", handler.Error)
		public.Post("/users", handler.CreateUser)
		public.Get("/feeds", handler.GetAllFeeds)
	})

	// private group
	v1ApiRouter.Group(func(private chi.Router) {
		private.Use(middleware.BasicAuth(dbConn))
		private.Get("/users", handler.GetMyUserData)
		private.Post("/feeds", handler.CreateFeed)
		private.Post("/feed_follows", handler.CreateFeedFollow)
		private.Get("/feed_follows", handler.GetFeedFollows)
		private.Delete("/feed_follows/{feedFollowID}", handler.DeleteFeedFollows)
		private.Get("/posts", handler.GetPostsForUser)
	})

	mainRouter.Mount("/v1", v1ApiRouter)

	svr := &http.Server{
		Handler: mainRouter,
		Addr:    ":" + port,
	}

	log.Printf("Start serving on: %v", svr.Addr)
	err = svr.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PORT:", port)
}
