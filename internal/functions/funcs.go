package functions

import (
	"encoding/json"
	"fmt"
	"gorcom/internal/models"
	"io/fs"
	"path/filepath"
)

// Walk возвращает слайс имён  файлов по маске what (с путями), исключая маску имени excluder
func Walk(what, excluder string) (filesToZip []string, err error) {

	// папка искомого файла
	folder := filepath.Dir(what)

	fmt.Println("On Unix:")

	err = filepath.Walk(folder,
		// path - перебираются имена (с путями) файлов в папке folder
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
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
				fmt.Printf("  filepath.Match  %v: %v\n", path, err)
				return err
			}
			// если совпадают имена искомого и текущего файла
			if matched {
				// то проверяем маску экслудера "exclude"
				matchedEx, err := filepath.Match(excluder, info.Name())
				if err != nil {
					fmt.Printf("  filepath.Match  %v: %v\n", path, err)
					return err
				}
				// если таки да, экслуде его
				if matchedEx {
					return nil
				}
				// если добрались сюда - значит всё совпадает, в массив его
				filesToZip = append(filesToZip, path)
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
