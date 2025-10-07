package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	middlewareecho "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	dbadapter "github.com/FIAP-11SOAT/totem-de-pedidos/internal/adapter/database"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/api"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/helper"
	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/middleware"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
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

	rdsSecretString, err := GetAwsSecrets("fase3-database-totem-de-pedidos-secrets")
	if err != nil {
		fmt.Println("Error to get RDS secret:", err)
		return
	}

	var rdsSecret RDSSecret
	err = json.Unmarshal([]byte(*rdsSecretString), &rdsSecret)
	if err != nil {
		fmt.Println("Error to parse RDS secret:", err)
		return
	}

	databaseAdapter := dbadapter.New(dbadapter.Input{
		DBDrive:   os.Getenv("DB_DRIVER"),
		DBUser:    rdsSecret.User,
		DBPass:    rdsSecret.Password,
		DBHost:    rdsSecret.Endpoint,
		DBName:    rdsSecret.DBName,
		DBOptions: os.Getenv("DB_OPTIONS"),
	})

	secretString, err := GetAwsSecrets("fase3-lambda-totem-de-pedidos-secrets")
	if err != nil {
		fmt.Println("Erro ao obter o segredo:", err)
		return
	}

	var cognitoSecret CognitoSecret
	err = json.Unmarshal([]byte(*secretString), &cognitoSecret)
	if err != nil {
		fmt.Println("Erro ao parsear o segredo:", err)
		return
	}

	keyMap, err := helper.ParseJWKS(cognitoSecret.CognitoJwksJson)
	if err != nil {
		log.Fatalf("Erro ao parsear JWKS: %v", err)
	}

	app := echo.New()
	app.Logger.SetLevel(log.INFO)

	app.Use(middlewareecho.CORS())
	app.Use(middlewareecho.Recover())

	app.Use(middleware.JWTAuthMiddleware(keyMap))

	api.Routers(app, databaseAdapter)

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

func GetAwsSecrets(secretName string) (*string, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar config AWS: %w", err)
	}
	client := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}
	output, err := client.GetSecretValue(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar secret: %w", err)
	}
	if output.SecretString == nil {
		return nil, errors.New("secret não encontrado ou vazio")
	}
	return output.SecretString, nil
}

type CognitoSecret struct {
	CognitoJwksJson string `json:"COGNITO_JWKS_JSON"`
}

type RDSSecret struct {
	Endpoint string `json:"RDS_ENDPOINT"`
	Port     int    `json:"RDS_PORT"`
	User     string `json:"RDS_USERNAME"`
	Password string `json:"RDS_PASSWORD"`
	DBName   string `json:"RDS_DATABASE"`
}
