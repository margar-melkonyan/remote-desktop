package router

import "github.com/go-chi/chi/v5"

func sessionsRouterGroup(sessions chi.Router) {
	sessions.Get("/", dependencies.SessionHandler.Get)
	sessions.Post("/", dependencies.SessionHandler.StoreConnection)
	sessions.Put("/{id}", dependencies.SessionHandler.UpdateConnection)
	sessions.Delete("/{id}", dependencies.SessionHandler.RemoveConnection)
	sessions.Get("/{id}/edit", dependencies.SessionHandler.Edit)
}
