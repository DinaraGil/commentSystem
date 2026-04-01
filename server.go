package main

import (
	"commentSystem/graph"
	"commentSystem/internal/models"
	"commentSystem/internal/storage"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

const defaultPort = "8080" //env
//env - in-memory or db

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

		store = storage.NewPostgresStorage(db)
		log.Println("connected to postgres")

	default:
		log.Fatalf("unknown STORAGE_TYPE: %s", storageType)
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			Store:                   store,
			CommentPublishedChannel: make(map[int][]chan *models.Comment),
		},
	}))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	//srv.SetQueryCache(lru.New)
	//srv.Use(extension.Introspection{})
	//srv.Use(extension.AutomaticPersistedQuery{
	//	Cache: lru.New,
	//})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("storage=%s", storageType)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
