package main

import (
	"commentSystem/graph"
	"commentSystem/graph/loaders"
	"commentSystem/internal/models"
	"commentSystem/internal/storage"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = defaultPort
	}

	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "memory"
	}

	var store storage.Storage
	var postgresDB *sqlx.DB

	switch storageType {
	case "memory":
		store = storage.NewMemoryStorage()
		log.Println("starting with in-memory storage")

	case "postgres":
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			log.Fatal("DATABASE_URL is required when STORAGE_TYPE=postgres")
		}

		db, err := sqlx.Connect("postgres", databaseURL)
		if err != nil {
			log.Fatalf("failed to connect to postgres: %v", err)
		}

		postgresDB = db
		store = storage.NewPostgresStorage(db)
		log.Println("connected to postgres")

	default:
		log.Fatalf("unknown STORAGE_TYPE: %s", storageType)
	}

	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					Store:                   store,
					CommentPublishedChannel: make(map[int]map[string]chan *models.Comment),
				},
			},
		),
	)

	srv.AddTransport(&transport.Websocket{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.Use(extension.Introspection{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	var gqlHandler http.Handler = srv

	if storageType == "postgres" && postgresDB != nil {
		gqlHandler = loaders.Middleware(postgresDB.DB, gqlHandler)
	}

	http.Handle("/query", gqlHandler)

	log.Printf("storage=%s", storageType)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
