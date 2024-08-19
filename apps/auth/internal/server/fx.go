package server

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Module() fx.Option {
	return fx.Module(
		"server",
		fx.Provide(NewConfig, NewOption, NewServer),
		fx.Invoke(
			func(lc fx.Lifecycle, server *grpc.Server, config Config, log *zap.Logger) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServerPort))
						if err != nil {
							return err
						}

						go func() {
							log.Sugar().Infof("Starting server on port %d", config.ServerPort)
							if err := server.Serve(lis); err != nil {
								log.Sugar().Errorf("failed to serve: %v", err)
							}
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						server.GracefulStop()
						return nil
					},
				})
			},
		),
	)
}