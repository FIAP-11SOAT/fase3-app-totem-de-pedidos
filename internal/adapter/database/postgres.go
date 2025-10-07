package dbadapter

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: implementar pool
// TODO: implementar abstração para transaction

type DatabaseAdapter struct {
	Client *pgxpool.Pool
}

type Input struct {
	DBDrive   string
	DBUser    string
	DBPass    string
	DBHost    string
	DBName    string
	DBOptions string
}

func New(input Input) *DatabaseAdapter {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	connStr := fmt.Sprintf("%s://%s:%s@%s/%s%s",
		input.DBDrive,
		input.DBUser,
		url.QueryEscape(input.DBPass),
		input.DBHost,
		input.DBName,
		input.DBOptions,
	)

	poolConfig := Config(connStr)

	client, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		panic(err)
	}

	if err = client.Ping(ctx); err != nil {

		slog.Error("error ping database", err.Error(), err)
		panic(err)
	}

	return &DatabaseAdapter{Client: client}
}

func (d *DatabaseAdapter) DataBaseHeatlh(ctx context.Context) error {
	return d.Client.Ping(ctx)
}

func Config(dbUrl string) *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	var DATABASE_URL string = dbUrl

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig
}
