package app

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/ehrktia/lendbook/internal/data"
	"github.com/ehrktia/lendbook/mocks"
	"github.com/jackc/pgx/v5/pgconn"
)

var ctx, cancel = context.WithCancel(context.Background())
var l = slog.Default()

func TestGetUserByEmailError(t *testing.T) {
	tests := []struct {
		name   string
		email  string
		ctxErr bool
		err    error
		expect func() *mocks.UserFetcher
	}{
		{
			name:  "no data error",
			email: "first.last@email.com",
			err:   data.ErrNodata,
			expect: func() *mocks.UserFetcher {
				mockUserDataStore := mocks.NewUserFetcher(t)
				mockUserDataStore.
					EXPECT().
					GetUserByEmail(ctx, "first.last@email.com").
					Return(int64(0), data.ErrNodata)
				return mockUserDataStore
			},
		},
		{
			name:   "ctx error",
			ctxErr: true,
			email:  "first.last@email.com",
			err:    ctx.Err(),
			expect: func() *mocks.UserFetcher {
				cancel()
				mockUserDataStore := mocks.NewUserFetcher(t)
				mockUserDataStore.
					EXPECT().GetUserByEmail(ctx, "first.last@email.com").
					Return(int64(0), ctx.Err())
				return mockUserDataStore
			},
		},
		{
			name:  "pg error",
			email: "first.last@email.com",
			err:   &pgconn.PgError{Code: t.Name(), Message: t.Name()},
			expect: func() *mocks.UserFetcher {
				cancel()
				mockUserDataStore := mocks.NewUserFetcher(t)
				mockUserDataStore.
					EXPECT().GetUserByEmail(ctx, "first.last@email.com").
					Return(int64(0), &pgconn.PgError{Code: t.Name(),
						Message: t.Name()})
				return mockUserDataStore
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserDataStore := test.expect()
			u := NewUser(mockUserDataStore, l)
			_, err := u.GetUserByEmail(ctx, test.email)
			if !mockUserDataStore.AssertExpectations(t) {
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
				if !errors.Is(err, test.err) {
					t.Fatalf("expected-%v,got-%v\n", test.err, err)
				}
			}

		})
	}
	t.Cleanup(func() { cancel() })
}
