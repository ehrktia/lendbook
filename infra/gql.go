package infra

import (
	"net/http"

	"codeberg.org/ehrktia/lendbook/internal/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func gqlServer(resolver *graph.Resolver, h *http.Server) error {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	mux := http.NewServeMux()
	addRoute(mux, srv)
	h.Handler = cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods: []string{http.MethodGet, http.MethodOptions,
			http.MethodConnect, http.MethodPost, http.MethodDelete,
			http.MethodPatch, http.MethodPut},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "token"},
	}).Handler(mux)
	h.Handler = optionsHandler(h.Handler)
	return nil
}
func addRoute(m *http.ServeMux, srv *handler.Server) {
	m.Handle("/", playground.Handler("GQL Playground", "/query"))
	m.Handle("/query", srv)
}

func optionsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, *")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
