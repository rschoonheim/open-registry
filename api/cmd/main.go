package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"open-registry-api/config"
	"open-registry-api/database"
	"open-registry-api/handlers"
	"open-registry-api/models"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to the database
	err = database.Connect(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	err = database.Migrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Create initial admin user if not exists
	createAdminUser()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg.JWTSecret)

	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Open Registry API")
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)

	// API routes (protected)
	api := app.Group("/api")
	api.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.JWTSecret),
	}))

	// Protected routes
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to protected API route",
		})
	})

	api.Get("/user", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		username := claims["username"].(string)

		return c.JSON(fiber.Map{
			"username": username,
			"admin":    claims["admin"],
		})
	})

	// Listen on configured port
	log.Printf("Starting server on port %s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}

// createAdminUser creates an admin user if it doesn't exist
func createAdminUser() {
	// Check if admin user exists
	var adminUser models.User
	result := database.DB.Where("username = ?", "admin").First(&adminUser)

	// If admin user doesn't exist, create one
	if result.Error != nil {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		// Create admin user
		adminUser = models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Email:    "admin@example.com",
			IsAdmin:  true,
		}

		result = database.DB.Create(&adminUser)
		if result.Error != nil {
			log.Fatalf("Failed to create admin user: %v", result.Error)
		}

		log.Println("Admin user created successfully")
	}
}
