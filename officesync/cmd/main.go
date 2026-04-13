package main

import (
	"log"
	"os"

	"officesync/handlers"
	"officesync/middleware"
	"officesync/models"
	"officesync/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Verify it exists in the root directory.")
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is not set.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to the OfficeSync database.")

	log.Println("Running database migrations...")
	err = db.AutoMigrate(&models.User{}, &models.Resource{}, &models.Booking{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	userService := services.NewUserService(db)
	resourceService := services.NewResourceService(db)
	bookingService := services.NewBookingService(db)

	userHandler := handlers.NewUserHandler(userService)
	resourceHandler := handlers.NewResourceHandler(resourceService)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	router := gin.Default()
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-User-Role", "X-User-Id"}
	router.Use(cors.New(corsConfig))

	router.Use(middleware.StructuredLogger())

	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", userHandler.HandleCreateUser)
		v1.GET("/resources", resourceHandler.HandleListResources)

		v1.POST("/resources", middleware.RequireAdmin(), resourceHandler.HandleCreateResource)

		v1.PUT("/resources/:id", middleware.RequireAdmin(), resourceHandler.HandleUpdateResource)

		v1.POST("/bookings", bookingHandler.HandleCreateBooking)
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
