//go:build integration

package app

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ehrktia/lendbook/internal/data"
	"github.com/ehrktia/lendbook/internal/data/pg"
)

func TestGetUserByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	once := &sync.Once{}
	conn, err := pg.NewPool(ctx, once)
	if err != nil {
		t.Fatal(err)
	}
	ur := data.NewUser(conn)
	id, err := ur.Create(ctx, data.User{
		FirstName: "fn" + t.Name(),
		LastName:  "ln-1" + t.Name(),
		Email:     "fn" + t.Name() + "." + "ln" + t.Name() + "@domain.com",
	})
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	t.Logf("user_id:%s\n", id)
	user, err := ur.GetById(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user:%#v\n", user)
	if id != user.ID {
		t.Fatalf("created user with id:%s\treceived user with id:%s\n", id, user.ID)
	}
	t.Cleanup(func() {
		if err := ur.DeleteByID(ctx, id); err != nil {
			t.Fatalf("err removing data:%v\n", err)
		}
		cancel()
	})
}
