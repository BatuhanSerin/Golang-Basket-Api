package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	check "github.com/BatuhanSerin/final-project/check"
	"github.com/BatuhanSerin/final-project/internal/User"
	"github.com/BatuhanSerin/final-project/internal/basket"
	"github.com/BatuhanSerin/final-project/internal/category"
	"github.com/BatuhanSerin/final-project/internal/product"
	"github.com/BatuhanSerin/final-project/package/config"
	db "github.com/BatuhanSerin/final-project/package/database"
	"github.com/BatuhanSerin/final-project/package/graceful"
	logger "github.com/BatuhanSerin/final-project/package/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	cfg, err := config.LoadConfig("./package/config/local-config")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	//Set global logger
	logger.NewLogger(cfg)
	defer logger.Close()

	//Set db
	DB := db.Connect(cfg)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerConfig.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.ServerConfig.ReadTimeoutSecs * int64(time.Second)),
		WriteTimeout: time.Duration(cfg.ServerConfig.WriteTimeoutSecs * int64(time.Second)),
	}

	//Set routes
	rootRouter := r.Group(cfg.ServerConfig.RoutePrefix)
	categoryRouter := rootRouter.Group("/categories")
	productRouter := rootRouter.Group("/products")
	userRouter := rootRouter.Group("/users")
	basketRouter := rootRouter.Group("/basket")

	// Product Repo
	productRepo := product.NewProductRepository(DB)
	productRepo.Migration()
	product.NewProductHandler(productRouter, productRepo, config.GetSecretKey())

	// Category Repo
	categoryRepo := category.NewCategoryRepository(DB)
	categoryRepo.Migration()
	category.NewCategoryHandler(categoryRouter, categoryRepo, config.GetSecretKey())

	// User Repo
	userRepo := User.NewUserRepository(DB)
	userRepo.Migration()
	User.NewUserHandler(userRouter, userRepo, config.GetSecretKey())

	// Basket Repo
	basketRepo := basket.NewBasketRepository(DB)
	basketRepo.Migration()
	basket.NewBasketHandler(basketRouter, basketRepo, config.GetSecretKey())

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("Listen error", zap.Error(err))
		}
	}()

	//CheckDependency checks database connection with ping and checks if server is running
	check.CheckDependency(DB, r)

	zap.L().Info("Server started", zap.String("port", cfg.ServerConfig.Port))

	graceful.ShutdownGin(srv, time.Duration(cfg.ServerConfig.TimeoutSecs*int64(time.Second)))
}
