package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-rtc-actions/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	port := flag.String("port", "8083", "default port is 8083")
	flag.Parse()
	addr := fmt.Sprintf(":%s", *port)

	router := gin.Default()

	router.GET("/ws", handler.ServeWS)

	server := http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    time.Duration(time.Minute * 3),
		WriteTimeout:   time.Duration(time.Minute * 3),
		MaxHeaderBytes: 10 << 20, //10MB
	}

	go func() {
		log.Println("server started listening on port : ", *port)
		if err := server.ListenAndServe(); err != nil {
			log.Println("server failed to initialized on port : ", *port)
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	// shutdown close
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("server failed to shutdown : %v\n", port)
	}
}
