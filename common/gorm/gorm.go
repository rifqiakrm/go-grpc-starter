package gorm

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresGormDB builds a connection of gorm to PostgreSQL.
func NewPostgresGormDB(pgxPool *pgxpool.Pool) (*gorm.DB, error) {
	pgxConn, err := pgxPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	conConfig := pgxConn.Conn().Config()
	conn := stdlib.OpenDB(*conConfig)

	connCfg, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	return connCfg, err
}
