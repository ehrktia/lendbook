package pg

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	connPool *pgxpool.Pool
}

func NewPool(ctx context.Context, once *sync.Once) (*Postgres, error) {
	var err error
	pg := &Postgres{}
	// TODO: move to env var or config
	connString := "postgres://postgres:postgres@localhost:5432/lendbook?"
	once.Do(func() {
		pg.connPool, err = pgxpool.New(ctx, connString)
	})
	if err != nil {
		return pg, err
	}
	return pg, nil
}

func (pg *Postgres) GetConn(ctx context.Context) (*pgxpool.Conn, error) {
	return pg.connPool.Acquire(ctx)
}

func RuninTx(
	ctx context.Context, conn *pgxpool.Conn, f func(tx pgx.Tx) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Conn().Close(ctx); err != nil {
			fmt.Printf("error closing conn:%v\n", err)
		}
		conn.Release()
	}()
	err = f(tx)
	if err != nil {
		rErr := tx.Rollback(ctx)
		return fmt.Errorf("%w:%v", rErr, err)
	}
	if err := tx.Commit(ctx); err != nil {
		fmt.Printf("error commiting txn:%v\n", err)
	}
	return nil
}
