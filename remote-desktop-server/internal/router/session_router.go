package router

import "github.com/go-chi/chi"

func sessionsRouterGroup(sessions chi.Router) {
	sessions.Get("/", dependencies.SessionHandler.Get)
	sessions.Get("/{id}/edit", dependencies.SessionHandler.Eidt)
	sessions.Post("/", dependencies.SessionHandler.StoreConnection)
	sessions.Put("/{id}", dependencies.SessionHandler.UpdateConnection)
	sessions.Delete("/{id}", dependencies.SessionHandler.RemoveConnection)
}
