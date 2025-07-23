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

	// чтение из JSON файла во втором агрументе CLI
	data, err := os.ReadFile(os.Args[2])
	if err != nil {
		models.Logger.Error("os.ReadFile  ", "err", err)
		return
	}

	// чтение конфигурации SSH
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

	if os.Args[1] == "create" {
		upa, err := functions.UnmarPack([]byte(data))
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", upa)
		err = functions.U0packer(upa)
		if err != nil {
			return err
		}
		return err
	}
	upa, err := functions.UnmarUnPack([]byte(data))
	functions.UnPack(upa)


	return
}
