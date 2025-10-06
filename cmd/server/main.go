package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	middlewareecho "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api"
)

func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {

	isProduction := getEnvOrDefault("PROFILE", "dev") == "prod"
	if !isProduction {
		pwd, _ := os.Getwd()
		envFile := fmt.Sprintf("%s/.env", pwd)
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	databaseAdapter := dbadapter.New(dbadapter.Input{
		DBDrive:   os.Getenv("DB_DRIVER"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASS"),
		DBHost:    os.Getenv("DB_HOST"),
		DBName:    os.Getenv("DB_NAME"),
		DBOptions: os.Getenv("DB_OPTIONS"),
	})

	app := echo.New()
	app.Logger.SetLevel(log.INFO)

	app.Use(middlewareecho.CORS())
	app.Use(middlewareecho.Recover())

	api.Routers(app, databaseAdapter)

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
