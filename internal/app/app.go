package app

import (
	"database/sql"
	"fmt"
	"golang-gin-app/internal/handlers"
	"golang-gin-app/internal/repository"
	"golang-gin-app/internal/service"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Config struct to hold configuration
type Config struct {
	Server struct {
		Port int
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
	DatabaseSecondary struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
	Log struct {
		Level  string
		Format string
	}
	JWT struct {
		Secret     string
		Expiration string
	}
}

type App struct {
	Router           *gin.Engine
	DB               *sql.DB
	DBSecondary      *sql.DB
	Config           *Config
	Service          *service.Service
	ServiceSecondary *service.Service
}

func NewApp() *App {
	config := loadConfig()
	db, err := initDB(config, "primary")
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to primary database: %v", err))
	}
	dbSecondary, err := initDB(config, "secondary")
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to secondary database: %v", err))
	}
	repo := repository.NewUserRepository(db)
	repoSecondary := repository.NewUserRepository(dbSecondary)
	svc := service.NewService(repo)
	svcSecondary := service.NewService(repoSecondary)
	router := gin.Default()
	app := &App{
		Router:           router,
		DB:               db,
		DBSecondary:      dbSecondary,
		Config:           config,
		Service:          svc,
		ServiceSecondary: svcSecondary,
	}
	app.initializeMiddleware()
	app.initializeRoutes()
	app.loadTemplates()
	return app
}

func loadConfig() *Config {
	// Load configuration from environment variables or use defaults
	return &Config{
		Server: struct {
			Port int
		}{Port: getEnvInt("SERVER_PORT", 8080)},
		Database: struct {
			Host     string
			Port     int
			User     string
			Password string
			Name     string
		}{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 3306),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "P@ssw0rd"),
			Name:     getEnv("DB_NAME", "dtxcasemgnt"),
		},
		DatabaseSecondary: struct {
			Host     string
			Port     int
			User     string
			Password string
			Name     string
		}{
			Host:     getEnv("DB_SECONDARY_HOST", "localhost"),
			Port:     getEnvInt("DB_SECONDARY_PORT", 3306),
			User:     getEnv("DB_SECONDARY_USER", "root"),
			Password: getEnv("DB_SECONDARY_PASSWORD", "P@ssw0rd"),
			Name:     getEnv("DB_SECONDARY_NAME", "another_database"),
		},
		Log: struct {
			Level  string
			Format string
		}{Level: "info", Format: "json"},
		JWT: struct {
			Secret     string
			Expiration string
		}{Secret: "your_jwt_secret", Expiration: "24h"},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func initDB(config *Config, dbType string) (*sql.DB, error) {
	var dsn string
	if dbType == "primary" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Database.User, config.Database.Password, config.Database.Host,
			config.Database.Port, config.Database.Name)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DatabaseSecondary.User, config.DatabaseSecondary.Password, config.DatabaseSecondary.Host,
			config.DatabaseSecondary.Port, config.DatabaseSecondary.Name)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (a *App) initializeRoutes() {
	// Initialize your routes here
	a.Router.GET("/hello", handlers.HelloHandler)
	a.Router.GET("/fake-users", handlers.GenerateFakeUsersFormHandler)
	a.Router.POST("/fake-users", handlers.GenerateFakeUsersHandler(a.Service))
	// Route for secondary database API
	a.Router.GET("/fake-users-secondary", handlers.GenerateFakeUsersFormHandler)
	a.Router.POST("/fake-users-secondary", handlers.GenerateFakeUsersHandler(a.ServiceSecondary))
}

func (a *App) initializeMiddleware() {
	// Initialize your middleware here
	a.Router.Use(gin.Logger())
	a.Router.Use(gin.Recovery())
}

func (a *App) loadTemplates() {
	// Load HTML templates
	a.Router.LoadHTMLGlob("templates/*")
}

func (a *App) Run(addr string) error {
	if addr == "" {
		addr = fmt.Sprintf(":%d", a.Config.Server.Port)
	}
	return a.Router.Run(addr)
}
