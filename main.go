package main

import (
	"io"
	"log"
	"net/http"
	"test/internal/handler"
	"test/internal/repository/cache"

	"github.com/bytedance/sonic"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func main() {
	cache, err := cache.New()
	if err != nil {
		log.Fatal(err)
	}
	// Create a new router & API
	router := chi.NewMux()

	// Huma config
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

	cfg := huma.DefaultConfig("Url Shortener", "0.0.1")
	cfg.FieldsOptionalByDefault = false
	cfg.Formats = map[string]huma.Format{"sonic": SonicJSON}
	cfg.DefaultFormat = "sonic"

	// Create server
	api := humachi.New(router, cfg)
	handlers := handler.New(cache)

	// Link API
	grp := huma.NewGroup(api, "/api/v1")

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

	// Start server
	http.ListenAndServe("127.0.0.1:8888", router)
}
