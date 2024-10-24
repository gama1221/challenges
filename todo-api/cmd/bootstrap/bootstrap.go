package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-api/internal/config"
	"todo-api/internal/db"
	"todo-api/internal/handler"
	TodoRepository "todo-api/internal/repo/todo/gocql_impl"
	UserRepository "todo-api/internal/repo/user/gocql_impl"
	"todo-api/internal/service"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func Setup() {
	//Create a new Gin instance
	engine := gin.Default()
	engine.Use(cors.Default())
	router := engine.Group("/api/v1")

	// Read the configuration.
	cfg := config.Get()

	dbConn, err := db.New(cfg)
	if err != nil {
		fmt.Println(err.Error())
	}
	opensearchURL := "http://localhost:9200"
	tr := TodoRepository.NewTodoRepository(dbConn.Session, opensearchURL)
	ur := UserRepository.NewUserRepository(dbConn.Session)

	ts := service.NewTodoService(tr, ur)
	us := service.NewUserService(ur)

	th := handler.NewTodoHandler(ts)
	uh := handler.NewUserHandler(us)

	th.RegisterTodoRoutes(router)
	uh.RegisterUserRoutes(router)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: engine,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server existing")
}
