//go:build wireinject
// +build wireinject

package di

import (
	wire "github.com/google/wire"
	"queue-manager/internal"
	"queue-manager/internal/controllers"
	"queue-manager/internal/providers"
	"queue-manager/internal/queues"
	"queue-manager/internal/structures"
	"queue-manager/internal/tcp"
)

func InitApp(cfg *structures.CliFlags) (*internal.App, error) {

	wire.Build(
		providers.NewConfigProvider,
		providers.NewLogProvider,
		queues.NewZstdCompressor,
		queues.NewQueueManager,
		queues.NewQueueIOManager,
		queues.NewScheduler,
		controllers.NewApiController,
		controllers.NewAdminController,
		tcp.NewTcpServer,
		internal.InitRoutes,
		internal.NewApp,
	)
	return nil, nil
}
