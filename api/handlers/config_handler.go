package handlers

import (
	"github.com/gofiber/fiber/v2"
	"open-registry-api/config"
)

// ConfigHandler handles configuration-related endpoints
type ConfigHandler struct {
	Config *config.Config
}

// NewConfigHandler creates a new config handler
func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{
		Config: cfg,
	}
}

// GetFeatures returns the enabled features from the configuration
func (h *ConfigHandler) GetFeatures(c *fiber.Ctx) error {
	// Build a features object based on the YAML configuration
	features := fiber.Map{
		"authentication": fiber.Map{
			"register": h.Config.YAML.Authentication.Register.Enabled,
		},
	}

	return c.JSON(fiber.Map{
		"features": features,
	})
}
