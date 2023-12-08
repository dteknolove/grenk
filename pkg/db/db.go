package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dteknolove/grenk/pkg/vip"
)

type DbService struct {
	DB  *pgxpool.Pool
	TX  pgx.Tx
	Err error
}

func NewDbService(ctx context.Context) *DbService {
	var errServDB error
	DB, errDB := New(ctx).DB()
	TX, errTX := New(ctx).TX()
	if errDB != nil {
		errServDB = errDB
	}
	if errTX != nil {
		errServDB = errTX
	}

	return &DbService{
		DB:  DB,
		TX:  TX,
		Err: errServDB,
	}
}

type PgConfig struct {
	Ctx context.Context
}

func New(ctx context.Context) *PgConfig {
	return &PgConfig{
		Ctx: ctx,
	}
}

var newErr error

func (p *PgConfig) TX() (pgx.Tx, error) {
	DB, err := p.DB()
	if err != nil {
		return nil, err
	}
	TX, errTX := DB.Begin(p.Ctx)
	if errTX != nil {
		return nil, errTX
	}

	return TX, nil
}

func (p *PgConfig) DB() (*pgxpool.Pool, error) {
	pool, errDB := p.Conn()
	if errDB != nil {
		newErr = errDB
	}
	acquire, err := pool.Acquire(p.Ctx)
	if err != nil {
		return nil, err
	}
	defer acquire.Release()

	return pool, newErr
}

func (p *PgConfig) Conn() (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(p.Ctx, 10*time.Second)
	defer cancel()

	dbURL, errURL := p.dbUrl()
	if errURL != nil {
		return nil, errURL
	}
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		newErr = err
	}

	config.MaxConns = int32(50)
	config.MaxConnLifetime = time.Duration(50)

	pool, errPool := pgxpool.NewWithConfig(ctx, config)
	if errPool != nil {
		newErr = errPool
	}
	_, errConnect := pool.Exec(ctx, ";")
	if errConnect != nil {
		newErr = errConnect
	}
	return pool, newErr
}

func (p *PgConfig) dbUrl() (string, error) {
	vipConf, errVip := vip.New().App()
	if errVip != nil {
		return "", errVip
	}

	dbUser := vipConf.DbUsername
	dbPass := vipConf.DbPassword
	dbHost := vipConf.DbHost
	dbPort := strconv.Itoa(vipConf.DbPort)
	dbName := vipConf.DbName
	dbSchema := vipConf.DbSchema
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?search_path=%s",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
		dbSchema)
	return url, nil
}
