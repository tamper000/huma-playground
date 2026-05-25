package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test/internal/handler"
	"test/internal/repository/cache"
	"time"

	"github.com/bytedance/sonic"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/doganarif/govisual"
	"github.com/go-chi/chi/v5"
)

const version = "0.0.1"

func main() {
	// Init temp DB
	cache, err := cache.New()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new router & API
	router := chi.NewMux()

	// Huma config
	cfg := humaConfig()

	// Create server
	api := humachi.New(router, cfg)
	handlers := handler.New(cache)
	registerHandlers(api, handlers)

	// Start server
	router2 := govisual.Wrap(
		router,
		govisual.WithRequestBodyLogging(true),
		govisual.WithResponseBodyLogging(true),
		govisual.WithIgnorePaths("/favicon.ico"),
		govisual.WithBasicAuth("admin", "admin"),
		govisual.WithLocalhostOnly(),
		govisual.WithMemoryStorage(),
	)
	server := &http.Server{
		Addr:    ":8888",
		Handler: router2,

		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 5,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}

func registerHandlers(api huma.API, handlers *handler.Handler) {
	// Link API
	grp := huma.NewGroup(api, "/api/v1")
	grp.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		ctx.SetHeader("Test-Header", "123")

		userAgent := ctx.Header("User-Agent")
		log.Println(userAgent)

		ctx = huma.WithValue(ctx, "userID", "123")
		next(ctx)
	})

	huma.Register(grp, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/shorten",
		Summary:     "Shorten link",
		Description: "Shorten link",
		Tags:        []string{"Links"},
		Errors:      []int{http.StatusInternalServerError},
	}, handlers.AddLink)

	huma.Register(grp, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/info/{id}",
		Summary:     "Info by ID",
		Description: "Getting info about a link by ID",
		Tags:        []string{"Links"},
		Errors:      []int{http.StatusNotFound},
	}, handlers.GetLink)

	// Redirect
	huma.Register(api, huma.Operation{
		Method:        http.MethodGet,
		Path:          "/{id}",
		Summary:       "Redirect to link",
		Description:   "Redirect to link by ID",
		Tags:          []string{"Redirect"},
		Errors:        []int{http.StatusNotFound},
		DefaultStatus: http.StatusTemporaryRedirect,
	}, handlers.RedirectLink)
}

func humaConfig() huma.Config {
	var SonicJSON = huma.Format{
		Marshal: func(writer io.Writer, v any) error {
			data, err := sonic.Marshal(v)
			if err != nil {
				return err
			}
			_, err = writer.Write(data)
			return err
		},
		Unmarshal: sonic.Unmarshal,
	}

	cfg := huma.DefaultConfig("Url Shortener", version)
	cfg.FieldsOptionalByDefault = false

	cfg.Formats = map[string]huma.Format{
		"application/json":         SonicJSON,
		"application/problem+json": SonicJSON,
	}

	cfg.DefaultFormat = "application/json"

	return cfg
}
