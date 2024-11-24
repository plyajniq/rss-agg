package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"rss-agg/internal/database"
	"rss-agg/internal/handler/api"
	"rss-agg/internal/handler/front"
	"rss-agg/internal/middleware"
	"rss-agg/internal/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "rss-agg/docs"

	_ "github.com/lib/pq"
)

//	@title			RSS AGG API
//	@version		1.0
//	@description	RSS aggregator with Chi in Go.
//	@host			localhost:8080
//	@BasePath		/api/v1
func main() {
	fmt.Println("Running...")

	// load env data
	godotenv.Load("./.env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is empty")
	}

	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal("HOST is empty")
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
	mainRouter.Use(middleware.DBConn(dbConn))

	mainRouter.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	mainRouter.Get("/swagger/*", httpSwagger.WrapHandler)

	apiRouter := chi.NewRouter()

	v1ApiRouter := chi.NewRouter()
	// public group
	v1ApiRouter.Group(func(public chi.Router) {
		public.Get("/healthz", api.HealthCheck)
		public.Get("/err", api.Error)
		public.Post("/users", api.CreateUser)
		public.Get("/feeds", api.GetAllFeeds)
	})

	// private group
	v1ApiRouter.Group(func(private chi.Router) {
		private.Use(middleware.BasicAuth(dbConn))
		private.Get("/users", api.GetMyUserData)
		private.Post("/feeds", api.CreateFeed)
		private.Post("/feed_follows", api.CreateFeedFollow)
		private.Get("/feed_follows", api.GetFeedFollows)
		private.Delete("/feed_follows/{feedFollowID}", api.DeleteFeedFollows)
		private.Get("/posts", api.GetPostsForUser)
	})

	apiRouter.Mount("/v1", v1ApiRouter)
	mainRouter.Mount("/api", apiRouter)

	basicWeb := chi.NewRouter()
	basicWeb.Group(func(public chi.Router) {
		public.Get("/", front.GetTopFeeds)
		public.Get("/feeds/{feedID}", front.GetFeedPosts)
		public.Get("/about", front.GetAbout)
	})

	mainRouter.Mount("/", basicWeb)

	svr := &http.Server{
		Handler: mainRouter,
		Addr:    host + ":" + port,
	}

	log.Printf("Start serving on: %v", svr.Addr)
	err = svr.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("HOST:", host)
	fmt.Println("PORT:", port)
}
