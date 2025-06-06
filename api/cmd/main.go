package main

import (
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
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Load environment variables from YAML config
	//
	err = godotenv.Load(cfg.YAML.Environment.File)
	if err != nil {
		log.Fatalf("Failed to load environment variables from YAML config: %v", err)
	}

	// Log whether registration is enabled based on YAML config
	if cfg.YAML.Authentication.Register.Enabled {
		log.Println("User registration is enabled")
	} else {
		log.Println("User registration is disabled")
	}

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
	authHandler := handlers.NewAuthHandler(cfg.Env.JWTSecret)
	configHandler := handlers.NewConfigHandler(cfg)

	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Open Registry API")
	})

	// Features endpoint - publicly accessible
	app.Get("/api/features", configHandler.GetFeatures)

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/login", authHandler.Login)

	// Only enable registration if it's enabled in the YAML config
	if cfg.YAML.Authentication.Register.Enabled {
		auth.Post("/register", authHandler.Register)
	}

	// API routes (protected)
	api := app.Group("/api")
	api.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.Env.JWTSecret),
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
	log.Printf("Starting server on port %s", cfg.Env.AppPort)
	log.Fatal(app.Listen(":" + cfg.Env.AppPort))
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
