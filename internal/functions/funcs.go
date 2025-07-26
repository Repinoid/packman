package functions

import (
	"encoding/json"
	"fmt"

	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"gorcom/internal/models"
)

type FileWinfo struct {
	FilePath string
	Info     os.FileInfo
}

// Walk возвращает слайс имён  файлов по маске what (с путями), исключая маску имени excluder
func Walk(what, excluder string) (filesToZip []FileWinfo, err error) {

	// папка искомого файла
	folder := filepath.Dir(what)

	models.Logger.Debug("what to pack %s, excluder is %s", what, excluder)

	err = filepath.Walk(folder,
		// path - перебираются имена (с путями) файлов в папке folder
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				models.Logger.Debug("prevent panic by handling failure accessing", "path", path, "err", err)
				return err
			}
			// если директория - пропускаем
			if info.IsDir() {
				return nil
			}
			// если папки искомого и текущего файла не совпадают возврат. отсекает и вложенные папки
			if folder != filepath.Dir(path) {
				return nil
			}
			matched, _ := filepath.Match("*ssh*", info.Name())
			if matched {
				return nil
			}
			// filepath.Base(what) - имя (или маска) искомого файла. info.Name() - имя текущего файла
			matched, err = filepath.Match(filepath.Base(what), info.Name())
			if err != nil {
				models.Logger.Error(" filepath.Match err ", "what", what, "err", err)
				return err
			}
			// если совпадают имена искомого и текущего файла
			if matched {
				// то проверяем маску экслудера "exclude"
				matchedEx, err := filepath.Match(excluder, info.Name())
				if err != nil {
					models.Logger.Error(" filepath.Match err ", "path", what, "err", err)
					return err
				}
				// если таки да, экслуде его
				if matchedEx {
					return nil
				}
				// если добрались сюда - значит всё совпадает, в массив его
				filesToZip = append(filesToZip, FileWinfo{FilePath: path, Info: info})
			}
			return nil
		})
	return
}

// UnmarPack анмаршаллит данные из файла с конфигурацией упаковки
func UnmarPack(data []byte) (u *models.Upack, err error) {
	upa := models.Upack{}
	err = json.Unmarshal([]byte(data), &upa)
	if err != nil {
		return
	}
	// выведем содержимое анмаршаленного 
	fmt.Printf("Pack name %s\nPack version %s\n", upa.Name, upa.Version)
	for i, v := range upa.Targets {
		fmt.Printf("Target %2d path %s ", i+1, v.(map[string]string)["path"])
		if excl, ok := v.(map[string]string)["exclude"]; ok {
			fmt.Printf("exclude %s", excl)
		}
		fmt.Println()
	}
	for i, v := range upa.Packets {
		fmt.Printf("Packet %2d name %s ", i+1, v.Name)
		if v.Ver != "" {
			fmt.Printf("version %s", v.Ver)
		}
		fmt.Println()
	}
	fmt.Println()
	return &upa, nil
}

// UnmarPack анмаршаллит данные из файла с конфигурацией распаковки
func UnmarUnPack(data []byte) (u models.Packages, err error) {
	upa := models.Packages{}
	err = json.Unmarshal(data, &upa)
	if err != nil {
		return
	}

	for i, v := range upa.Packs {
		fmt.Printf("Packet %2d name %s ", i+1, v.Name)
		if v.Ver != "" {
			fmt.Printf("version %s", v.Ver)
		}
		fmt.Println()
	}
	fmt.Println()


	return upa, nil
}

// ParseComparisonWithRegex определяет корректность условия по версии пакета, возвращает строку операции сравнения и строку с номером версии
func ParseComparisonWithRegex(expr string) (op, right string, err error) {
	// Регулярное выражение для операторов сравнения
	re := regexp.MustCompile(`^\s*(.*?)\s*(>=|<=|==|!=|>|<)\s*(.*?)\s*$`)

	matches := re.FindStringSubmatch(expr)
	if len(matches) != 4 {
		return "", "", fmt.Errorf("invalid comparison expression")
	}
	re = regexp.MustCompile(`^[0-9.]+$`)
	matched := re.MatchString(matches[3])
	if !matched {
		return "", "", fmt.Errorf("invalid comparison expression")
	}

	return matches[2], matches[3], nil
}

// compara применяет логику операции сравнения, заданную в строке (типа ">=")
func compara(left, op, right string) bool {
	verOk := false
	switch op {
	case "==":
		if left == right {
			verOk = true
		}
	case "!=":
		if left != right {
			verOk = true
		}
	case "<":
		if left < right {
			verOk = true
		}
	case ">":
		if left > right {
			verOk = true
		}
	case ">=":
		if left >= right {
			verOk = true
		}
	case "<=":
		if left <= right {
			verOk = true
		}
	}
	return verOk
}

// https://go.dev/play/p/j5B0nr55_or
