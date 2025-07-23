package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"gorcom/internal/functions"
	"gorcom/internal/models"
	"log/slog"
	"os"
	"strings"
)

func main() {

	ctx := context.Background()

	if len(os.Args) < 2 || (strings.ToLower(os.Args[1]) != "create" && strings.ToLower(os.Args[1]) != "update") {
		fmt.Println("Неверные аргументы. Запуск программы pm create ./packet.json или pm update ./packages.json")
		os.Exit(1)
	}

	// Если есть флаг -debug
	Level := slog.LevelInfo
	isDebug := false
	restoreFlag := flag.Bool("debug", isDebug, "Минимальный уровень логирования")
	flag.Parse()
	if *restoreFlag {
		Level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     Level,
		AddSource: true, // Добавлять информацию об исходном коде
	})
	models.Logger = slog.New(handler)
	slog.SetDefault(models.Logger)

	if err := Run(ctx); err != nil {
		models.Logger.Error(err.Error())
	}

}

func Run(ctx context.Context) (err error) {

	data, err := os.ReadFile(os.Args[2])
	if err != nil {
		models.Logger.Error("os.ReadFile  ", "err", err)
		return
	}

	sshConfData, err := os.ReadFile("sshConf.json")
	if err != nil {
		models.Logger.Error("Ошибка чтения конфигурации файлa sshConf.json ", "err", err)
		return
	}

	err = json.Unmarshal([]byte(sshConfData), &models.SSHConf)
	if err != nil {
		models.Logger.Error("Ошибка в конфигурации SSH, файл sshConf.json ", "err", err)
		return
	}

	upa, err := functions.UnmarPack([]byte(data))
	if err != nil {
		return
	}

	fmt.Printf("%+v\n", upa)

	err = functions.U0packer(upa)

	// op, right, err := functions.ParseComparisonWithRegex("  <  2.0  ")
	// fmt.Printf("Parsed: op='%s', right='%s' %v\n", op, right, err)

	return
}
