//go:build integration

package app

import (
	"fmt"
	"strconv"
	"testing"

	"codeberg.org/ehrktia/lendbook/internal/data"
)

func TestGetUserByID(t *testing.T) {
	id, err := userRepo.Create(ctx, data.User{
		FirstName: "fn" + t.Name(),
		LastName:  "ln-1" + t.Name(),
		Email:     "fn" + t.Name() + "." + "ln" + t.Name() + "@domain.com",
	})
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	t.Logf("user_id:%s\n", id)
	user, err := userRepo.GetById(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user:%#v\n", user)
	if id != user.ID {
		t.Fatalf("created user with id:%s\treceived user with id:%s\n", id, user.ID)
	}
	t.Cleanup(func() {
		if err := userRepo.DeleteByID(ctx, id); err != nil {
			t.Fatalf("err removing data:%v\n", err)
		}
	})
}

func TestGetAllPrevNext(t *testing.T) {
	of, limit := "0", "3"
	result, err := bookApp.GetAll(ctx, of, limit)
	if err != nil {
		t.Fatal(err)
	}
	if result.Prev != of {
		t.Fatalf("prev value - expected:%s\tgot:%s\n", of, result.Prev)
	}
	n, err := strconv.Atoi(limit)
	if err != nil {
		t.Fatal(err)
	}
	next := n + 1
	if result.Next != fmt.Sprintf("%d", next) {
		t.Fatalf("next value - expected:%s\tgot:%s\n",
			fmt.Sprintf("%d", next),
			result.Next)
	}
	t.Cleanup(func() { cancel() })
}
