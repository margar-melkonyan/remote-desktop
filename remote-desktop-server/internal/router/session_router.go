package router

import "github.com/go-chi/chi"

func sessionsRouterGroup(sessions chi.Router) {
	sessions.Get("/", dependencies.SessionHandler.Get)
	sessions.Post("/add-connection", dependencies.SessionHandler.AddNewConnection)
}
