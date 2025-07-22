package main

import (
	"context"
	"fmt"
	"gorcom/internal/functions"
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

	data := `{
				"name": "packet-1",
				"ver": "1.10",
				"targets": [
					"./archive_this1/*.txt",
					{"path": "./archive_this2/*", "exclude": "*.tmp"}
						],
				"packets": [
				{"name": "packet-3", "ver": "<=2.0"}
				]
				}`

	_ = data
	upa, err := functions.Unmar([]byte(data))
	if err != nil {
		return
	}

	fmt.Printf("%+v\n", upa)

	a, err := functions.Walk("../../cmd/packman/*pack*", "*.go*")
	fmt.Printf("%+v %v\n", a, err)

	return nil
}
