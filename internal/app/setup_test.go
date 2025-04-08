package app

import (
	"context"
	"log/slog"
	"sync"

	"codeberg.org/ehrktia/lendbook/internal/data"
	"codeberg.org/ehrktia/lendbook/internal/data/pg"
)

var (
	once          = &sync.Once{}
	ctx, cancel   = context.WithCancel(context.Background())
	connPool, err = pg.NewPool(ctx, once)
	l             = slog.Default()
	userRepo      = data.NewUser(connPool)
	bookRepo      = data.NewBook(connPool)
	bookApp       = NewBook(bookRepo, l)
)
