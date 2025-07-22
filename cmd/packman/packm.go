package main

import (
	"context"
	"gorcom/internal/models"
	"log/slog"
	"os"
)

func main() {

	ctx := context.Background()

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug, // Минимальный уровень логирования
		AddSource: true,            // Добавлять информацию об исходном коде
	})
	models.Logger = slog.New(handler)
	slog.SetDefault(models.Logger)

	if err := Run(ctx); err != nil {
		models.Logger.Error(err.Error())
	}

}

func Run(ctx context.Context) (err error) {

	

	return
}
