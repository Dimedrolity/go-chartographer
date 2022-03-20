package app

import (
	"chartographer-go/chart"
	"chartographer-go/imagetile"
	"chartographer-go/pkg/kvstore"
	"chartographer-go/server"
)

// Run инициализирует зависимости сервера и запускает его.
func Run(port, dataDirPath string, tileMaxSize int) error {
	tileRepo, err := imagetile.NewFileSystemTileRepo(dataDirPath)
	if err != nil {
		return err
	}
	bmpService := imagetile.NewBmpService(tileRepo)
	imageRepo := kvstore.NewInMemoryStore()
	chartService := chart.NewChartographerService(imageRepo, bmpService, tileMaxSize)
	config := server.NewConfig(port)
	srv := server.NewServer(config, chartService)

	return srv.Run()
}
