package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

// User represents the structure of our user
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Set of demo users
var users = map[string]string{
	"admin": "admin123",
	"user":  "user123",
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found")
	}

	// Set JWT secret from environment variables or use default
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-256-bit-secret"
		fmt.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	// Initialize Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Open Registry API")
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/login", login(jwtSecret))

	// API routes (protected)
	api := app.Group("/api")
	api.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
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

	// Example packages route
	api.Get("/packages", getPackages)

	// Listen on port from env or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}

// Login handler
func login(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		// Check if user exists and password is correct
		password, exists := users[user.Username]
		if !exists || password != user.Password {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		// Create the Claims
		claims := jwt.MapClaims{
			"username": user.Username,
			"admin":    user.Username == "admin",
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token
		t, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		return c.JSON(fiber.Map{
			"token": t,
			"user": fiber.Map{
				"username": user.Username,
				"admin":    user.Username == "admin",
			},
		})
	}
}

// Example of a handler for getting packages
func getPackages(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"packages": []fiber.Map{
			{
				"id":          "1",
				"name":        "example-package",
				"version":     "1.0.0",
				"description": "An example package",
			},
			{
				"id":          "2",
				"name":        "another-package",
				"version":     "2.1.0",
				"description": "Another example package",
			},
		},
	})
}
