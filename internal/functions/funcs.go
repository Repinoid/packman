package functions

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"gorcom/internal/models"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type FileWinfo struct {
	FilePath string
	Info     os.FileInfo
}

// Walk возвращает слайс имён  файлов по маске what (с путями), исключая маску имени excluder
func Walk(what, excluder string) (filesToZip []FileWinfo, err error) {

	// папка искомого файла
	folder := filepath.Dir(what)

	fmt.Println("On Unix:")

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
			// filepath.Base(what) - имя (или маска) искомого файла. info.Name() - имя текущего файла
			matched, err := filepath.Match(filepath.Base(what), info.Name())
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

func Unmar(data []byte) (u *models.Upack, err error) {
	upa := models.Upack{}
	err = json.Unmarshal([]byte(data), &upa)
	if err != nil {
		return
	}
	return &upa, nil
}

func Upacker(upa *models.Upack) (err error) {
	for _, u := range upa.Targets {
		exclude := u.(map[string]string)["exclude"]
		filesToZip, err := Walk(u.(map[string]string)["path"], exclude)
		if err != nil {
			return err
		}
		fmt.Printf("%+v %v\n", filesToZip, err)

		op, right := "", ""
		for _, pack := range upa.Packets {
			if pack.Ver != "" {
				op, right, err = ParseComparisonWithRegex(pack.Ver)
				if err != nil {
					return fmt.Errorf("ошибка условия версии  %s: %v", pack.Ver, err)
				}
			}
			// если не выполляется условие, заданное в версии пакета, пример  {"name": "packet-3", "ver": "<="2.0" },
			if !compara(upa.Version, op, right) {
				continue
			}

			zf, err := os.Create(pack.Name)
			if err != nil {
				return fmt.Errorf("ошибка создания файла %s: %v", pack.Name, err)
			}
			defer zf.Close()

			zipWriter := zip.NewWriter(zf)
			defer zipWriter.Close()

			// Устанавливаем комментарий с версией
			zipWriter.SetComment(fmt.Sprintf("Version: %s", upa.Version))

			for _, f := range filesToZip {
				header, err := zip.FileInfoHeader(f.Info)
				if err != nil {
					return err
				}
				writer, err := zipWriter.CreateHeader(header)
				if err != nil {
					return err
				}
				fileToArchive, err := os.Open(f.FilePath)
				if err != nil {
					return err
				}
				defer fileToArchive.Close()

				_, err = io.Copy(writer, fileToArchive)
				if err != nil {
					return err
				}
			}
		}

	}

	return
}

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

// https://go.dev/play/p/j5B0nr55_or

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
