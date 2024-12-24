//go:build integration

package pg

import (
	"context"
	"sync"
	"testing"
)

func TestNewPool(t *testing.T) {
	ctx := context.Background()
	once := &sync.Once{}
	var err error
	var p *Postgres
	t.Run("create new conn pool", func(t *testing.T) {
		p, err = NewPool(ctx, once)
		if err != nil {
			t.Fatal(err)
		}
		if p == nil {
			t.Fatalf("expected postgres,got nil")
		}
	})
	t.Run("get conn", func(t *testing.T) {
		conn, err := p.GetConn(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if conn == nil {
			t.Fatalf("expected conn from pool,got nil")
		}
		t.Cleanup(func() {
			if err := conn.Conn().Close(ctx); err != nil {
				t.Fatal(err)
			}
			conn.Release()
		})
	})
}
