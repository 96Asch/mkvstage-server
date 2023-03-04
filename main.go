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

func setupMigrations(gormDatabase *gorm.DB) error {
	models := [...]any{
		&domain.User{},
		&domain.Bundle{},
		&domain.Song{},
		&domain.Role{},
		&domain.UserRole{},
	}

	for _, model := range models {
		log.Printf("Inserting table %s", reflect.TypeOf(model))

		err := gormDatabase.AutoMigrate(model)
		if err != nil {
			return domain.NewInitializationErr(err.Error())
		}
	}

	return nil
}

func setupStore() (*gorm.DB, *redis.Client) {
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_NAME")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")

	var gormDatabase *gorm.DB

	var err error

	gormDatabase, err = store.GetDB(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatal(err)
	}

	err = setupMigrations(gormDatabase)
	if err != nil {
		log.Fatal(err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	rdb, err := store.GetRedis(redisHost, redisPort)
	if err != nil {
		log.Fatal(err)
	}

	return gormDatabase, rdb
}

func main() {
	router := gin.Default()
	database, tokenDatabase := setupStore()

	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")

	userRepo := repository.NewGormUserRepository(database)
	tokenRepo := repository.NewRedisTokenRepository(tokenDatabase)
	bundleRepo := repository.NewGormBundleRepository(database)
	songRepo := repository.NewGormSongRepository(database)
	userroleRepo := repository.NewGormUserRoleRepository(database)
	roleRepo := repository.NewGormRoleRepository(database)
	setlistRepo := repository.NewGormSetlistRepository(database)

	userService := service.NewUserService(userRepo, roleRepo, userroleRepo)
	tokenService := service.NewTokenService(tokenRepo, userRepo, accessSecret, refreshSecret)
	mhw := middleware.NewGinMiddlewareHandler(tokenService)
	bundleService := service.NewBundleService(bundleRepo)
	songService := service.NewSongService(userRepo, songRepo)
	userroleService := service.NewUserRoleService(userroleRepo)
	roleService := service.NewRoleService(roleRepo, userRepo, userroleRepo)
	setlistService := service.NewSetlistService(userRepo, setlistRepo)

	config := handler.Config{
		Router: router,
		U:      userService,
		T:      tokenService,
		MH:     mhw,
		B:      bundleService,
		S:      songService,
		R:      roleService,
		UR:     userroleService,
		SL:     setlistService,
	}

	run(&config)
}
