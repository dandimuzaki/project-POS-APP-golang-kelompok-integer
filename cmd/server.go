package cmd

import (
	"fmt"
	"log"
	"net/http"
	"travel-api/internal/wire"
)

func APiserver(app *wire.App) {
	fmt.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", app.Route); err != nil {
		log.Fatal("can't run service")
	}

	// // gracefully shutdown ------------------------------------------------------------------------
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit

	// close(app.Stop)
	// app.WG.Wait()

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Println("can't shutdown service")
	// }

	// log.Println("server shutdown cleanly")
}

