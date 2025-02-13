package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"

	gohttp "net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hafiztri123/internal/adapters/primary/http"
	"github.com/hafiztri123/internal/adapters/secondary/persistent"
	"github.com/hafiztri123/internal/core/entity"
	"github.com/hafiztri123/internal/core/usecase"
	"github.com/hafiztri123/pkg/config"
	"github.com/hafiztri123/pkg/middleware"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const BASE_URL = "/api/v1"

func main() {
	setupGoDotEnv()
	cfg := config.LoadConfig()
	db := DatabaseInit(cfg)
	createTable(db)
	router := gin.Default()
	corsConfig := corsInit(cfg)
	router.Use(cors.New(*corsConfig))
	useCors(router, corsConfig)
	ApplyProductionSetting(router, cfg)
	jwt := JWTInit(cfg)


	authHandler := AuthHandlerInit(db, jwt)
	taskHandler := TaskHandlerInit(db)
	categoryHandler := CategoryHandlerInit(db)



	AuthRoutesInit(router, authHandler)
	TaskRoutesInit(router, taskHandler, jwt)
	CategoryRoutesInit(router, categoryHandler, jwt)

	StartServer(cfg, router)




}



func DatabaseInit(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error
	maxRetries := 5
	var logType logger.LogLevel

	if cfg.Environment == "production"{
		logType = logger.Silent
	} else {
		logType = logger.Info
	}

 

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
			Logger: logger.Default.LogMode(logType),
		})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database, attempt %d/%d", i+1, maxRetries)
		time.Sleep(time.Second * 5)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after multiple attempts:", err)
	}

	return db

}

func createTable(db *gorm.DB) {
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Task{})
}

func AuthHandlerInit(db *gorm.DB, jwt *usecase.JWTService) *http.AuthHandler {

	repo := persistent.NewAuthRepository(db)
	service :=  usecase.NewAuthService(repo)
	handler := http.NewAuthHandler(service, jwt)
	return handler
}

func AuthRoutesInit(router *gin.Engine, handler *http.AuthHandler) {
	auth := router.Group(BASE_URL + "/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}

}

func TaskRoutesInit(router *gin.Engine, handler *http.TaskHandler, jwt *usecase.JWTService) {
	protected := router.Group(BASE_URL + "/tasks") 
	protected.Use(middleware.AuthMiddleware(jwt))
	{
		protected.POST("", handler.CreateTask)
		protected.GET("", handler.GetUserTasks)
		protected.GET("/:task_id/subtasks", handler.GetSubTasks)
		protected.PUT("/:task_id", handler.UpdateTask)
		protected.DELETE("/:task_id", handler.DeleteTask)
		protected.GET("/category/:category_id", handler.GetTaskByCategory)
		protected.GET("/non-category", handler.GetNonCategoryTasks)

	}
	

}

func CategoryRoutesInit(router *gin.Engine, handler *http.CategoryHandler, jwt *usecase.JWTService) {
	protected := router.Group(BASE_URL + "/category")
	protected.Use(middleware.AuthMiddleware(jwt))
	{
		protected.POST("", handler.CreateCategory)
		protected.GET("", handler.GetAllCategory)
		protected.DELETE("/:category_id", handler.DeleteCategory)
	}
}





func setupGoDotEnv(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env")
	}
}


func JWTInit(cfg *config.Config) *usecase.JWTService {
	jwtconfig := &config.JWTConfig{
		SecretKey: cfg.JWTSecret,
		TokenLifetime: int(cfg.JWTExpiration.Minutes()),
	}

	jwtService := usecase.NewJWTService(jwtconfig)
	return jwtService
}

func TaskHandlerInit(db *gorm.DB) *http.TaskHandler {

	taskRepo := persistent.NewTaskRepository(db)
	userRepo := persistent.NewUserRepository(db)

	service := usecase.NewTaskService(taskRepo, userRepo)
	handler := http.NewTaskHandler(service)
	return handler
}

func CategoryHandlerInit(db *gorm.DB) *http.CategoryHandler {

	categoryRepo := persistent.NewCategoryRepository(db)
	service := usecase.NewCategoryService(categoryRepo)
	handler := http.NewCategoryHandler(service)
	return handler
}

func corsInit(cfg *config.Config) *cors.Config{
	return &cors.Config{
		AllowOrigins:     cfg.AllowedOrigins, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}
}

func useCors (router *gin.Engine, corsConfig *cors.Config) {
	router.Use(cors.New(*corsConfig))
}

func ApplyProductionSetting(router *gin.Engine, cfg *config.Config) {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
		router.Use(middleware.RateLimit(100, time.Minute))
		router.Use(middleware.SecurityHeaders())
	}
}

func StartServer(cfg *config.Config, router *gin.Engine ) {
	srv := &gohttp.Server{
		Addr: ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != gohttp.ErrServerClosed {
			log.Fatalf("Failed to start server:%v" ,err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err :=srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}


}