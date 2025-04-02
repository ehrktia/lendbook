//go:build integration

package data

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"codeberg.org/ehrktia/lendbook/internal/data/pg"
)

var ctx = context.Background()
var once = &sync.Once{}
var lenderId string
var pool *pg.Postgres

func TestMain(m *testing.M) {
	if err := setupTestData(); err != nil {
		fmt.Printf("setup failed: %v\n", err)
		os.Exit(1)
	}
	exitCode := m.Run()
	if err := cleanUp(); err != nil {
		fmt.Printf("error cleaning up: %v\n", err)
	}
	os.Exit(exitCode)
}

func TestUserById(t *testing.T) {
	tCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	lender := NewUser(pool)
	id, err := lender.Create(ctx, User{
		FirstName: t.Name(),
		LastName:  t.Name(),
		Email:     "first.last@domain.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	o, err := lender.GetById(tCtx, id)
	if err != nil {
		t.Fatal(err)
	}
	if id != o.ID {
		t.Fatalf("expected:%s\tgot:%s\n", id, o.ID)

	}
}

var queryInsBook = `INSERT INTO public.book(
 title, author, edition, owner_id, available, added, updated)
	VALUES ( 'book-1', 'author-1', '1', $1, TRUE, now(), now());`

var queryInsUser = `INSERT INTO public.lender(
	first_name, last_name, email, active,version)
	VALUES ('first', 'last', 'first.last@email.com', true,1);`

var queryUserByfirst = `select id from "lender" where first_name='first';`

func setupTestData() error {
	var err error
	pool, err = pg.NewPool(ctx, once)
	if err != nil {
		return err
	}

	conn, err := pool.GetConn(ctx)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	tag, err := tx.Exec(ctx, queryInsUser)
	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return errors.New("expected to insert one row failed")
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	_ = conn.QueryRow(ctx, queryUserByfirst).Scan(&lenderId)

	tx, err = conn.Begin(ctx)
	if err != nil {
		return err
	}

	tag, err = tx.Exec(ctx, queryInsBook, lenderId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() < 1 {
		return errors.New("expected to insert one row failed")
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	if err := tx.Conn().Close(ctx); err != nil {
		return err
	}
	conn.Release()
	return nil
}

func TestCreateUser(t *testing.T) {
	tCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var err error
	lender := NewUser(pool)
	got, err := lender.Create(tCtx,
		User{
			FirstName: t.Name(), LastName: t.Name(),
			Email: t.Name()})
	if err != nil {
		t.Fatal(err)
	}
	if got == "" {
		t.Fatalf("expected valid uuid\tgot:%s\n", got)
	}

}

func TestGetUsers(t *testing.T) {
	tCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	lender := NewUser(pool)
	lenders, err := lender.GetUsers(tCtx)
	if err != nil {
		t.Fatal(err)
	}
	if len(lenders) < 2 {
		t.Fatalf("expected:2,got:%d", len(lenders))
	}
}

func TestUpdateUser(t *testing.T) {
	tCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	lender := NewUser(pool)
	id, err := lender.GetUserByEmail(tCtx, `first.last@email.com`)
	if err != nil {
		t.Fatal(err)
	}
	testDataList := []UserWithNoBooks{
		{
			ID:        id,
			FirstName: "karthick-1",
			LastName:  "lastname-1",
			Email:     "karthie-1.lastname-1@email.com",
			Active:    true,
			Version:   2,
		},
		{
			ID:        id,
			FirstName: "karthick-2",
			LastName:  "lastname-2",
			Email:     "karthie-2.lastname-2@email.com",
			Active:    true,
			Version:   3,
		},
	}
	ow, err := lender.Update(tCtx, testDataList[0])
	if err != nil {
		t.Fatal(err)
	}
	if ow.Version != testDataList[0].Version {
		t.Fatalf("expected:%s,got:%s\n", testDataList[0].String(), ow.String())
	}
	ow, err = lender.Update(tCtx, testDataList[1])
	if err != nil {
		t.Fatal(err)
	}
	if ow.Version != testDataList[1].Version {
		t.Fatalf("expected:%s,got:%s\n", testDataList[1].String(), ow.String())
	}

}

func cleanUp() error {
	conn, err := pool.GetConn(ctx)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `truncate table public.lender cascade;`); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `truncate table public.book cascade`); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil

}
