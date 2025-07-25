package functions

import (
	"archive/zip"
	"fmt"
	"sync"

	"io"
	"os"

	"gorcom/internal/models"
	"gorcom/internal/ssher"
)

func U0packer(upa *models.Upack) (err error) {

	var wg sync.WaitGroup
	wg.Add(len(upa.Packets))
	ch := make(chan error)

	// итерация по пакетам - определение в какой архив писать и ограничение по версиям
	for _, pack := range upa.Packets {
		// op - строка операции сравнения версии, right - значение версии
		op, right := "", ""
		// если для пакета версия "ver" задана  - {"name": "packet-3", "ver": "<="2.0" },
		if pack.Ver != "" {
			// парсим строку версии
			op, right, err = ParseComparisonWithRegex(pack.Ver)
			if err != nil {
				return fmt.Errorf("ошибка условия версии  %s: %v", pack.Ver, err)
			}
			// если не выполняется условие, заданное в версии пакета, пример  {"name": "packet-3", "ver": "<="2.0" },
			if !compara(upa.Version, op, right) {
				continue
			}
		}
		// анонимизация дабы не плодить сущности. упаковка в архив и засыл по SSH горутиной
		err = func(pName string) (err error) {
			zf, err := os.Create(pName)
			if err != nil {
				return fmt.Errorf("ошибка создания файла %s: %v", pName, err)
			}
			defer zf.Close()

			zipWriter := zip.NewWriter(zf)
			defer zipWriter.Close() // !!!! чортов defer

			// итерация по папкам заданным в "targets"
			for _, targa := range upa.Targets {
				// exclude, ok  - если ключа "exclude" нет, exclude = ""
				exclude, ok := targa.(map[string]string)["exclude"]
				_ = ok
				// возвращает слайс имён  файлов для упаковки
				filesToZip, err := Walk(targa.(map[string]string)["path"], exclude)
				if err != nil {
					return err
				}
				// Устанавливаем комментарий с версией
				zipWriter.SetComment(upa.Version)
				// пакуем файлы, прошедшие отбор, в архив
				for _, f := range filesToZip {
					header, err := zip.FileInfoHeader(f.Info)
					if err != nil {
						return err
					}
					header.Method = zip.Deflate
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
			go func(ch chan<- error) {
				defer wg.Done()
				err = ssher.LoadBySSH(models.SSHConf.Host, models.SSHConf.User, models.SSHConf.Password, pName, "/files/"+pName)
				ch <- err
			}(ch)
			return
		}(pack.Name)

	}

	go func() {
		wg.Wait()
		close(ch)
	}()
	// Читаем результаты из канала
	for result := range ch {
		fmt.Println(result)
	}

	return
}
