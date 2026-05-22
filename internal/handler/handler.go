package handler

import (
	"context"
	"net/http"
	"test/internal/repository/cache"
	"test/internal/utils"
)

type Handler struct {
	cache *cache.Cache
}

func New(cache *cache.Cache) *Handler {
	return &Handler{
		cache: cache,
	}
}

func (h *Handler) AddLink(ctx context.Context, input *CreateShortLink) (*ShortLinkOutput, error) {
	url := input.Body.URL

	ID, err := utils.RandomID()
	if err != nil {
		return nil, unknownError
	}

	err = h.cache.CreateLink(ctx, url, ID)
	if err != nil {
		return nil, unknownError
	}

	resp := &ShortLinkOutput{}
	resp.Body.ID = ID
	return resp, nil
}

func (h *Handler) GetLink(ctx context.Context, input *GetLink) (*InfoLink, error) {
	id := input.ID

	link, err := h.cache.GetLink(ctx, id)
	if err != nil {
		return nil, IDNotExists
	}

	info := &InfoLink{}
	info.Body.Link = link

	return info, nil
}

func (h *Handler) RedirectLink(ctx context.Context, input *GetLink) (*RedirectOutput, error) {
	id := input.ID

	link, err := h.cache.GetLink(ctx, id)
	if err != nil {
		return nil, IDNotExists
	}

	return &RedirectOutput{
		Status:   http.StatusTemporaryRedirect,
		Location: link,
	}, nil
}
