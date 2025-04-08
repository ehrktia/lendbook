package app

import (
	"context"
	"errors"
	"strings"
	"testing"

	"codeberg.org/ehrktia/lendbook/internal/data"
	"codeberg.org/ehrktia/lendbook/mocks"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestGetUserByEmailError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	tests := []struct {
		name   string
		email  string
		ctxErr bool
		err    error
		expect func() *mocks.UserQuery
	}{
		{
			name:  "no data error",
			email: "first.last@email.com",
			err:   data.ErrNodata,
			expect: func() *mocks.UserQuery {
				mockUserDataStore := mocks.NewUserQuery(t)
				mockUserDataStore.
					EXPECT().
					GetUserByEmail(ctx, "first.last@email.com").
					Return("", data.ErrNodata)
				return mockUserDataStore
			},
		},
		{
			name:   "ctx error",
			ctxErr: true,
			email:  "first.last@email.com",
			err:    ctx.Err(),
			expect: func() *mocks.UserQuery {
				cancel()
				mockUserDataStore := mocks.NewUserQuery(t)
				mockUserDataStore.
					EXPECT().GetUserByEmail(ctx, "first.last@email.com").
					Return("", ctx.Err())
				return mockUserDataStore
			},
		},
		{
			name:  "pg error",
			email: "first.last@email.com",
			err:   &pgconn.PgError{Code: t.Name(), Message: t.Name()},
			expect: func() *mocks.UserQuery {
				cancel()
				mockUserDataStore := mocks.NewUserQuery(t)
				mockUserDataStore.
					EXPECT().GetUserByEmail(ctx, "first.last@email.com").
					Return("",
						&pgconn.PgError{Code: t.Name(), Message: t.Name()})
				return mockUserDataStore
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mocksQuerier := test.expect()
			mocksCommand := mocks.NewCommander(t)
			u := NewUser(mocksCommand, mocksQuerier, l)
			_, err := u.GetUserByEmail(ctx, test.email)
			if !mocksQuerier.AssertExpectations(t) {
				t.Fatal("expected calls not made")
			}
			if err == nil {
				t.Fatalf("expected:%v,got:%v\n", ctx.Err(), err)
			}
			switch {
			case test.ctxErr:
				if !errors.Is(err, ctx.Err()) {
					t.Fatalf("expected:%v,got:%v\n", ctx.Err(), err)
				}
			default:
				if !strings.EqualFold(err.Error(), test.err.Error()) {
					t.Fatalf("expected-%v,got-%v\n", test.err, err)
				}
			}

		})
	}
	t.Cleanup(func() { cancel() })
}
