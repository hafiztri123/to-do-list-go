package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"

	"github.com/gin-gonic/gin"
	"github.com/hafiztri123/internal/adapters/primary/http"
	"github.com/hafiztri123/internal/adapters/secondary/persistent"
	"github.com/hafiztri123/internal/core/entity"
	"github.com/hafiztri123/internal/core/usecase"
	"github.com/hafiztri123/pkg/config"
	"github.com/hafiztri123/pkg/logger"
	"github.com/hafiztri123/pkg/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"github.com/gin-contrib/cors"
)

const BASE_URL = "/api/v1"

func main() {
	setupZeroLog()
	setupGoDotEnv()
	dbconfig := &config.DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "hafizh",
		Password: "Sudarmi12",
		DBname:   "test_db",
	}

	gormDB := connectGormToDB(dbconfig)
	createTable(gormDB)

	jwt := setupJWT()
	authHandler := createAuthHandler(gormDB, jwt)
	taskHandler := createTaskHandler(gormDB)
	categoryHandler := createCategoryHandler(gormDB)


	router := gin.Default()
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())
	corsSetting := setupCors()
	router.Use(cors.New(corsSetting))

	setupAuthRoutes(router, authHandler)
	setupTaskRoutes(router, taskHandler, jwt)
	setupCategoryRoutes(router, categoryHandler, jwt)

	placeholderRoutes(router, jwt)

	router.Run(":8080")


}



func connectGormToDB(config *config.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Host, config.User, config.Password, config.DBname, config.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func createTable(db *gorm.DB) {
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Task{})
}

func createAuthHandler(db *gorm.DB, jwt *usecase.JWTService) *http.AuthHandler {

	repo := persistent.NewAuthRepository(db)
	service :=  usecase.NewAuthService(repo)
	handler := http.NewAuthHandler(service, jwt)
	return handler
}

func setupAuthRoutes(router *gin.Engine, handler *http.AuthHandler) {
	auth := router.Group(BASE_URL + "/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}

}

func setupTaskRoutes(router *gin.Engine, handler *http.TaskHandler, jwt *usecase.JWTService) {
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

func setupCategoryRoutes(router *gin.Engine, handler *http.CategoryHandler, jwt *usecase.JWTService) {
	protected := router.Group(BASE_URL + "/category")
	protected.Use(middleware.AuthMiddleware(jwt))
	{
		protected.POST("", handler.CreateCategory)
		protected.GET("", handler.GetAllCategory)
		protected.DELETE("/:category_id", handler.DeleteCategory)
	}
}

func placeholderRoutes(router *gin.Engine, jwt *usecase.JWTService) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(jwt))
	protected.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Protected route accessed successfully",
		})
	})
}

func setupZeroLog(){
	logger.Init()

	log.Info().Str("service", "auth").Msg("Starting auth service")
}

func setupGoDotEnv(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
}


func setupJWT() *usecase.JWTService {
	jwtconfig := &config.JWTConfig{
		SecretKey: os.Getenv("JWT_SECRET"),
		TokenLifetime: 3600,
	}

	jwtService := usecase.NewJWTService(jwtconfig)
	return jwtService
}

func createTaskHandler(db *gorm.DB) *http.TaskHandler {

	taskRepo := persistent.NewTaskRepository(db)
	userRepo := persistent.NewUserRepository(db)

	service := usecase.NewTaskService(taskRepo, userRepo)
	handler := http.NewTaskHandler(service)
	return handler
}

func createCategoryHandler(db *gorm.DB) *http.CategoryHandler {

	categoryRepo := persistent.NewCategoryRepository(db)
	service := usecase.NewCategoryService(categoryRepo)
	handler := http.NewCategoryHandler(service)
	return handler
}

func setupCors() cors.Config{
	return cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
}