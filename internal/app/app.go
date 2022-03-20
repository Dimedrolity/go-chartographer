package app

import (
	"go-chartographer/internal/chart"
	imagetile2 "go-chartographer/internal/imagetile"
	"go-chartographer/internal/server"
	"go-chartographer/pkg/kvstore"
)

// Run инициализирует зависимости сервера и запускает его.
func Run(port, dataDirPath string, tileMaxSize int) error {
	tileRepo, err := imagetile2.NewFileSystemTileRepo(dataDirPath)
	if err != nil {
		return err
	}
	bmpService := imagetile2.NewBmpService(tileRepo)
	imageRepo := kvstore.NewInMemoryStore()
	chartService := chart.NewChartographerService(imageRepo, bmpService, tileMaxSize)
	config := server.NewConfig(port)
	srv := server.NewServer(config, chartService)

	return srv.Run()
}
