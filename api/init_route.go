package api

import (
	"github.com/go-chi/chi/v5"
)

func InitRoute(mux chi.Router) {
	InitWorkflowRoute(mux)
}
