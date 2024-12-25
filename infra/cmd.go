package infra

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/ehrktia/lendbook/internal/app"
	"github.com/ehrktia/lendbook/internal/data"
	"github.com/ehrktia/lendbook/internal/data/pg"
	"github.com/ehrktia/lendbook/internal/graph"
)

func Run(ctx context.Context, l *slog.Logger) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	pgSync := &sync.Once{}
	postgres, err := pg.NewPool(ctx, pgSync)
	if err != nil {
		return err
	}
	userRepo := data.NewUser(postgres)
	// bookRepo := data.NewBook(postgres)
	userApp := app.NewUser(userRepo, l)
	resolver := &graph.Resolver{
		UserService: userApp,
	}
	port := getGQLServerPort()
	httpServer := initWebServer(port)
	if err := gqlServer(resolver, httpServer); err != nil {
		return err
	}
	l.LogAttrs(ctx, slog.LevelInfo,
		"starting gql server", slog.String("port", port))
	go func() {
		<-interrupt
		l.LogAttrs(ctx, slog.LevelInfo, "closing gql server")
		if err := httpServer.Shutdown(ctx); err != nil {
			l.LogAttrs(ctx, slog.LevelError, "error closing gql server",
				slog.String("error", err.Error()))
			os.Exit(1)
		}

	}()
	return httpServer.ListenAndServe()
}

func getGQLServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	return port

}

func initWebServer(port string) *http.Server {
	return &http.Server{
		Addr: "[::]:" + port,
	}

}
