package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/oktavarium/gophkeeper/internal/client/grpcclient"
	"github.com/oktavarium/gophkeeper/internal/server/grpcserver"
	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func TestClientServer(t *testing.T) {
	type setupFields struct {
		clientStorage *grpcclient.Mockstorage
		serverStorage *grpcserver.MockStorage
	}

	tests := []struct {
		name    string
		setup   func(f *setupFields)
		call    func(ctx context.Context, c *grpcclient.GrpcClient) error
		wantErr bool
	}{
		{
			name: "ok registration",
			setup: func(f *setupFields) {
				gomock.InOrder(
					f.clientStorage.EXPECT().GetToken(),
					f.serverStorage.EXPECT().Register(gomock.Any(), "user", "pass"),
					f.serverStorage.EXPECT().GetUserIDByLogin(gomock.Any(), "user"),
					f.serverStorage.EXPECT().UpdateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
					f.clientStorage.EXPECT().GetToken(),
					f.clientStorage.EXPECT().UpdateToken(gomock.Any(), gomock.Any()),
				)
			},
			call: func(ctx context.Context, c *grpcclient.GrpcClient) error {
				return c.Register(ctx, models.UserInfo{
					Login:    "user",
					Password: "pass",
				})
			},
			wantErr: false,
		},
		{
			name: "bad registration with same user",
			setup: func(f *setupFields) {
				gomock.InOrder(
					f.clientStorage.EXPECT().GetToken(),
					f.serverStorage.EXPECT().Register(gomock.Any(), "user", "pass").Return(fmt.Errorf("user already exists")),
				)
			},
			call: func(ctx context.Context, c *grpcclient.GrpcClient) error {
				return c.Register(ctx, models.UserInfo{
					Login:    "user",
					Password: "pass",
				})
			},
			wantErr: true,
		},
		{
			name: "ok login",
			setup: func(f *setupFields) {
				gomock.InOrder(
					f.clientStorage.EXPECT().GetToken(),
					f.serverStorage.EXPECT().GetUserIDByLogin(gomock.Any(), "user"),
					f.serverStorage.EXPECT().Login(gomock.Any(), "user", "pass"),
					f.serverStorage.EXPECT().UpdateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
					f.clientStorage.EXPECT().GetToken(),
					f.clientStorage.EXPECT().UpdateToken(gomock.Any(), gomock.Any()),
				)
			},
			call: func(ctx context.Context, c *grpcclient.GrpcClient) error {
				return c.Login(ctx, models.UserInfo{
					Login:    "user",
					Password: "pass",
				})
			},
			wantErr: false,
		},
		{
			name: "sync",
			setup: func(f *setupFields) {
				gomock.InOrder(
					f.clientStorage.EXPECT().GetDataEncrypted(),
					f.clientStorage.EXPECT().GetToken().Return("token_id", time.Now().Add(5*time.Minute), nil),
					f.serverStorage.EXPECT().GetToken(gomock.Any(), "token_id").Return("127.0.0.1", time.Now().Add(5*time.Minute), nil),
					f.serverStorage.EXPECT().GetUserIDByToken(gomock.Any(), "token_id"),
					f.serverStorage.EXPECT().GetUserIDByToken(gomock.Any(), "token_id").Return("user_id", nil),
					f.serverStorage.EXPECT().Sync(gomock.Any(), "user_id", gomock.Any()),
					f.serverStorage.EXPECT().UpdateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
					f.clientStorage.EXPECT().GetToken(),
					f.clientStorage.EXPECT().UpdateToken(gomock.Any(), gomock.Any()),
				)
			},
			call: func(ctx context.Context, c *grpcclient.GrpcClient) error {
				return c.Sync(ctx)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := setupFields{
				clientStorage: grpcclient.NewMockstorage(ctrl),
				serverStorage: grpcserver.NewMockStorage(ctrl),
			}

			ctx := context.Background()
			s, err := grpcserver.NewGrpcServer(ctx, f.serverStorage, "localhost:8888", "rootCACert.pem", "rootCAKey.pem")
			require.NoError(t, err)
			eg := new(errgroup.Group)
			eg.Go(func() error {
				return s.ListenAndServe()
			})

			time.Sleep(100 * time.Millisecond)

			c, err := grpcclient.NewGrpcClient(f.clientStorage)
			require.NoError(t, err)
			err = c.Init(ctx, "localhost:8888")
			require.NoError(t, err)

			if tt.setup != nil {
				tt.setup(&f)
			}
			if err := tt.call(ctx, c); (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			s.Stop()
			eg.Wait()
		})
	}
}
