package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/appleboy/graceful"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	var envfile string
	flag.StringVar(&envfile, "env-file", ".env", "Read in a file of environment variables")
	flag.Parse()

	_ = godotenv.Load(envfile)
	cfg, err := Environ()
	if err != nil {
		log.Fatal("invalid configuration")
	}

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(requestid.New())
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	app.GET("/healthz", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})
	app.NoRoute(func(c *gin.Context) {
	})

	m := graceful.NewManager()

	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           app.Handler(),
		ReadHeaderTimeout: time.Minute,
		ReadTimeout:       0,
		WriteTimeout:      0,
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
	}

	m.AddRunningJob(func(ctx context.Context) error {
		log.Printf("server running on %s port", cfg.Server.Port)
		return listenAndServe(srv)
	})
	m.AddShutdownJob(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	})

	<-m.Done()
}

func listenAndServe(s *http.Server) error {
	return s.ListenAndServe()
}
