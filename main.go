package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/96Asch/mkvstage-server/internal/handler"
	"github.com/96Asch/mkvstage-server/internal/handler/middleware"
	"github.com/96Asch/mkvstage-server/internal/repository"
	"github.com/96Asch/mkvstage-server/internal/service"
	"github.com/96Asch/mkvstage-server/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func run(config *handler.Config) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	handler.Initialize(config)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: config.Router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func setupMigrations(db *gorm.DB) {

	domains := [...]any{
		&domain.User{},
		&domain.Bundle{},
		&domain.Song{},
	}

	for _, domain := range domains {
		log.Printf("Inserting table %s", reflect.TypeOf(domain))
		db.AutoMigrate(domain)
	}

}

func setupStore() (*gorm.DB, *redis.Client) {
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_NAME")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")

	var db *gorm.DB
	var err error
	db, err = store.GetDB(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatal(err)
	}

	setupMigrations(db)

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	var rdb *redis.Client
	rdb, err = store.GetRedis(redisHost, redisPort)
	if err != nil {
		log.Fatal(err)
	}

	return db, rdb
}

func main() {

	router := gin.Default()
	db, rdb := setupStore()

	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")

	ur := repository.NewGormUserRepository(db)
	us := service.NewUserService(ur)

	tr := repository.NewRedisTokenRepository(rdb)
	ts := service.NewTokenService(tr, ur, accessSecret, refreshSecret)

	mhw := middleware.NewGinMiddlewareHandler(ts)

	br := repository.NewGormBundleRepository(db)
	bs := service.NewBundleService(br)

	sr := repository.NewGormSongRepository(db)
	ss := service.NewSongService(ur, sr)

	config := handler.Config{
		Router: router,
		U:      us,
		T:      ts,
		MH:     mhw,
		B:      bs,
		S:      ss,
	}

	run(&config)
}
