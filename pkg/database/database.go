package database

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/bariscan97/clean-rest-architecture/pkg/config"
)

func NewConnection(cfg *config.Config) *pgxpool.Pool {
    dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.DBName,
        cfg.Database.SSLMode,
    )

    poolConfig, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        panic(err)
    }

    pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
    if err != nil {
        panic(err)
    }

    return pool
}
