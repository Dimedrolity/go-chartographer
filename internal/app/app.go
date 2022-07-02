package app

import (
	"github.com/Dimedrolity/go-chartographer/internal/chart"
	"github.com/Dimedrolity/go-chartographer/internal/imgstore"
	"github.com/Dimedrolity/go-chartographer/internal/server"
	"github.com/Dimedrolity/go-chartographer/pkg/kvstore"
)

// Run инициализирует зависимости сервера и запускает его.
func Run(port, dataDirPath string, tileMaxSize int) error {
	tileRepo, err := imgstore.NewFileSystemTileRepo(dataDirPath)
	if err != nil {
		return err
	}

	bmpService := imgstore.NewBmpService(tileRepo)

	imageRepo := kvstore.NewInMemoryStore()
	adapter := &chart.ImageAdapter{}
	chartService := chart.NewChartographerService(imageRepo, bmpService, adapter, tileMaxSize)
	config := server.NewConfig(port)
	srv := server.NewServer(config, chartService)

	return srv.Run()
}
