package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lehoanglong/whatdish/internal/engagement"
	"github.com/lehoanglong/whatdish/internal/nutrition"
	"github.com/lehoanglong/whatdish/internal/platform/notion"
	"github.com/lehoanglong/whatdish/internal/platform/supabase"
	"github.com/lehoanglong/whatdish/internal/recipe"
	"github.com/lehoanglong/whatdish/internal/shared/config"
	"github.com/lehoanglong/whatdish/internal/shared/database"
	"github.com/lehoanglong/whatdish/internal/shared/middleware"
	"github.com/lehoanglong/whatdish/internal/suggestion"
	"github.com/lehoanglong/whatdish/internal/user"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Setup structured logging
	setupLogger(cfg.LogLevel)

	slog.Info("starting WhatDish API",
		"environment", cfg.Environment,
		"port", cfg.Port,
	)

	// Database connection
	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Platform adapters
	authClient := supabase.NewAuthClient(cfg.SupabaseURL, cfg.SupabaseAnonKey, cfg.SupabaseServiceRoleKey)
	notionClient := notion.NewClient(cfg.NotionAPIKey)

	// --- Recipe context ---
	dishRepo := recipe.NewDishRepository(pool)
	catRepo := recipe.NewCategoryRepository(pool)
	tagRepo := recipe.NewTagRepository(pool)
	dishService := recipe.NewDishService(dishRepo, catRepo, tagRepo)
	dishHandler := recipe.NewDishHandler(dishService)
	syncService := recipe.NewSyncService(notionClient, cfg.NotionDatabaseID, dishRepo, catRepo, tagRepo)

	// --- Suggestion context ---
	sessionRepo := suggestion.NewSessionRepo(pool)
	configRepo := suggestion.NewConfigRepo(pool)
	exclusionRepo := suggestion.NewExclusionRepo(pool)
	dishReader := suggestion.NewDishReaderAdapter(pool)
	calorieProvider := suggestion.NewCalorieProviderAdapter(pool)

	strategies := map[string]suggestion.Strategy{
		"RANDOM":      suggestion.NewRandomStrategy(dishReader, exclusionRepo, sessionRepo),
		"BY_CALORIES": suggestion.NewCalorieStrategy(dishReader, calorieProvider, sessionRepo),
		"BY_GROUP":    suggestion.NewGroupStrategy(dishReader, calorieProvider, sessionRepo),
	}

	suggestionService := suggestion.NewSuggestionService(strategies, sessionRepo, configRepo)
	suggestionHandler := suggestion.NewSuggestionHandler(suggestionService)

	// --- User context ---
	profileRepo := user.NewProfileRepo(pool)
	allergyRepo := user.NewAllergyRepo(pool)
	tdeeCalc := nutrition.NewTDEECalculator()

	authService := user.NewAuthService(authClient)
	profileService := user.NewProfileService(profileRepo, allergyRepo, tdeeCalc)
	userHandler := user.NewHandler(authService, profileService)

	// --- Nutrition context ---
	nutritionRepo := nutrition.NewNutritionRepo(pool)
	goalRepo := nutrition.NewGoalRepo(pool)
	nutritionService := nutrition.NewNutritionService(nutritionRepo, goalRepo)
	nutritionHandler := nutrition.NewHandler(nutritionService)

	// --- Engagement context ---
	favoriteRepo := engagement.NewFavoriteRepo(pool)
	viewRepo := engagement.NewViewRepo(pool)
	ratingRepo := engagement.NewRatingRepo(pool)
	engagementService := engagement.NewEngagementService(favoriteRepo, viewRepo, ratingRepo)
	engagementHandler := engagement.NewHandler(engagementService)

	// Setup Gin router
	if !cfg.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(
		middleware.Recovery(),
		middleware.Logger(),
		middleware.CORS(cfg.FrontendURL),
	)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")

	// Public: Recipes
	v1.GET("/recipes", dishHandler.HandleListDishes)
	v1.GET("/recipes/random", dishHandler.HandleGetRandomDish)
	v1.GET("/recipes/search", dishHandler.HandleSearchDishes)
	v1.GET("/recipes/:id", dishHandler.HandleGetDish)
	v1.GET("/categories", dishHandler.HandleListCategories)
	v1.GET("/tags", dishHandler.HandleListTags)

	// Public: Suggestions
	v1.POST("/suggestions/random", suggestionHandler.HandleRandomSuggestion)
	v1.POST("/suggestions/by-calories", suggestionHandler.HandleCalorieSuggestion)
	v1.POST("/suggestions/by-group", suggestionHandler.HandleGroupSuggestion)

	// Public: Auth
	v1.POST("/auth/register", userHandler.HandleRegister)
	v1.POST("/auth/login", userHandler.HandleLogin)
	v1.POST("/auth/refresh", userHandler.HandleRefresh)

	// Public: Nutrition
	v1.GET("/nutrition/recipes/:id", nutritionHandler.HandleGetRecipeNutrition)
	v1.GET("/nutrition/recipes", nutritionHandler.HandleGetBatchNutrition)
	v1.POST("/nutrition/calculate-tdee", nutritionHandler.HandleCalculateTDEE)
	v1.GET("/nutrition/goals", nutritionHandler.HandleListGoals)

	// Public: Views
	v1.POST("/views", engagementHandler.HandleRecordView)

	// Protected: require auth
	protected := v1.Group("")
	protected.Use(middleware.RequireAuth(authClient))
	{
		protected.GET("/users/me", userHandler.HandleGetProfile)
		protected.PUT("/users/me/profile", userHandler.HandleUpdateProfile)

		protected.GET("/suggestions/history", suggestionHandler.HandleGetHistory)

		protected.POST("/favorites", engagementHandler.HandleAddFavorite)
		protected.DELETE("/favorites/:recipe_id", engagementHandler.HandleRemoveFavorite)
		protected.GET("/favorites", engagementHandler.HandleListFavorites)
		protected.GET("/favorites/check", engagementHandler.HandleCheckFavorites)
	}

	// Admin endpoints
	admin := v1.Group("/admin")
	admin.Use(middleware.RequireAuth(authClient))
	admin.Use(middleware.RequireAdmin(cfg.ParseAdminUserIDs()))
	{
		admin.POST("/sync", func(c *gin.Context) {
			result, err := syncService.Sync(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "sync failed",
					"message": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": result})
		})

		admin.POST("/nutrition/recipes/:id", nutritionHandler.HandleUpsertNutrition)
	}

	// Start sync scheduler (background)
	if cfg.NotionAPIKey != "" && cfg.NotionDatabaseID != "" {
		go startSyncScheduler(ctx, syncService, cfg.SyncIntervalMinutes)
	}

	// Start server
	srv := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("server started", "addr", cfg.Addr())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server stopped")
}

func setupLogger(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	slog.SetDefault(slog.New(handler))
}

func startSyncScheduler(ctx context.Context, syncService *recipe.SyncService, intervalMinutes int) {
	// Run initial sync on startup
	slog.Info("running initial Notion sync...")
	if result, err := syncService.Sync(ctx); err != nil {
		slog.Error("initial sync failed", "error", err)
	} else {
		slog.Info("initial sync completed", "added", result.Added, "updated", result.Updated)
	}

	// Schedule periodic sync
	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			slog.Info("running scheduled Notion sync...")
			if result, err := syncService.Sync(ctx); err != nil {
				slog.Error("scheduled sync failed", "error", err)
			} else {
				slog.Info("scheduled sync completed", "added", result.Added, "updated", result.Updated)
			}
		}
	}
}
