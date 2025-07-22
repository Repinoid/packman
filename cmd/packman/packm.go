package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gorcom/internal/models"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
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
	upa, err := unmar([]byte(data))
	if err != nil {
		return
	}

	fmt.Printf("%+v\n", upa)

	Walk()

	return nil
}

func unmar(data []byte) (u *models.Upack, err error) {
	upa := models.Upack{}
	err = json.Unmarshal([]byte(data), &upa)
	if err != nil {
		return
	}
	return &upa, nil
}

func Walk() {

	fmt.Println("On Unix:")
	err := filepath.Walk("../", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return nil
		}
		matched, err := filepath.Match("*packa*", info.Name())
		if err != nil {
			fmt.Printf("  filepath.Match  %v: %v\n", path, err)
			return err
		}

		fmt.Printf("%v  visited file or dir: %v\n", path, matched)
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", "../", err)
		return
	}
}
