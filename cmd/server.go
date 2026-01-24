package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project-POS-APP-golang-integer/internal/wire"
	"syscall"
	"time"
)

func APiserver(app *wire.App) {
	fmt.Printf("Server running on port %d", app.Config.Port)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.Config.Port),
		Handler: app.Route,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("can't run service")
		}
	}()

	// gracefully shutdown ------------------------------------------------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	close(app.Stop)
	app.WG.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("can't shutdown service")
	}

	log.Println("server shutdown cleanly")
}

