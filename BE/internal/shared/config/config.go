package config

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port        int    `envconfig:"PORT" default:"8080"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`

	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`

	SupabaseURL            string `envconfig:"SUPABASE_URL"`
	SupabaseAnonKey        string `envconfig:"SUPABASE_ANON_KEY"`
	SupabaseServiceRoleKey string `envconfig:"SUPABASE_SERVICE_ROLE_KEY"`

	NotionAPIKey     string `envconfig:"NOTION_API_KEY"`
	NotionDatabaseID string `envconfig:"NOTION_DATABASE_ID"`

	SyncIntervalMinutes int `envconfig:"SYNC_INTERVAL_MINUTES" default:"30"`

	FrontendURL string `envconfig:"FRONTEND_URL" default:"http://localhost:3000"`

	AdminUserIDs string `envconfig:"ADMIN_USER_IDS" default:""`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}

// ParseAdminUserIDs parses the comma-separated ADMIN_USER_IDS into a slice of UUIDs.
func (c *Config) ParseAdminUserIDs() []uuid.UUID {
	if c.AdminUserIDs == "" {
		return nil
	}
	parts := strings.Split(c.AdminUserIDs, ",")
	ids := make([]uuid.UUID, 0, len(parts))
	for _, p := range parts {
		id, err := uuid.Parse(strings.TrimSpace(p))
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
