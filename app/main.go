package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func newRouter() *httprouter.Router {
	mux := httprouter.New()
	godotenv.Load()

	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	if ytApiKey == "" {
		log.Fatal("error: youtube api key not set, termintating...")
	}
	ytChannelId := os.Getenv("YOUTUBE_CHANNEL_ID")
	if ytChannelId == "" {
		log.Fatal("error: channel id not set, termintating...")
	}

	mux.GET("/youtube/channel/stats", getChannelStats(ytApiKey, ytChannelId))

	return mux
}

func main() {
	srv := &http.Server{
		Addr:    ":10101",
		Handler: newRouter(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		log.Println("Service interrupt received")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error: %v", err)
		}

		log.Print("shutdown complete")
		close(idleConnsClosed)

	}()

	log.Printf("starting server on port 10101")
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fatal https server failed to start: %v", err)
		}
	}

	<-idleConnsClosed
	log.Println("service stop")
}
