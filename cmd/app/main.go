package main

import (
	"fmt"
	"github.com/bertpersyn/posology-graphql/internal/posology"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bertpersyn/posology-graphql/internal/graphql/graph"
	"github.com/bertpersyn/posology-graphql/internal/graphql/graph/generated"
	samparser "github.com/bertpersyn/posology-graphql/internal/sam"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	_ "net/http/pprof"
)

var version, gitCommit, application string

const defaultPort = "8080"

//todo: add metrics
//todo: add config
//todo: move init code away to service

func main() {
	_, err := fmt.Println(fmt.Sprintf("application: %v, gitCommit: %v, version: %v", application, gitCommit, version))
	if err != nil {
	}
	go func() {
		logrus.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	logrus.SetLevel(logrus.DebugLevel)
	samParserService, err := samparser.New()
	if err != nil {
		panic(err)
	}
	err = samParserService.ParseAll()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := &graph.Resolver{Sam: samParserService, Posology: posology.New()}
	r.Init()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	srv := c.Handler(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: r})))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	logrus.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	logrus.Fatal(http.ListenAndServe(":"+port, nil))
}
