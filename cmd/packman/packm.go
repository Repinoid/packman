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

	// data, err := os.ReadFile("packet.json")
	// if err != nil {
	// 	models.Logger.Error("os.ReadFile  ", "err", err)
	// 	return
	// }

	// upa, err := functions.Unmar([]byte(data))
	// if err != nil {
	// 	return
	// }

	// fmt.Printf("%+v\n", upa)

	// err = functions.Upacker(upa)

	op, right, err := functions.ParseComparisonWithRegex("  <  2.0  ")
	fmt.Printf("Parsed: op='%s', right='%s' %v\n", op, right, err)

	return
}
