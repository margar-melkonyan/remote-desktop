package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/margar-melkonyan/remote-desktop.git/internal/common/dependency"
	"github.com/margar-melkonyan/remote-desktop.git/internal/config"
	"github.com/margar-melkonyan/remote-desktop.git/internal/router"
)

func RunHttpServer() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()
	go func() {
		addr := fmt.Sprintf(":%s", config.ServerConfig.Port)
		server := &http.Server{
			Addr:    addr,
			Handler: router.NewRouter(dependency.NewAppDependencies()),
		}
		slog.Info(
			fmt.Sprintf("Http Server start on port %s",
				addr,
			),
		)
		err := server.ListenAndServe()
		if err != nil {
			return
		}
	}()
	<-ctx.Done()
}
