package functions

import (
	"archive/zip"
	"fmt"

	"io"
	"os"

	"gorcom/internal/models"
	"gorcom/internal/ssher"
)

func U0packer(upa *models.Upack) (err error) {

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

		zf, err := os.Create(pack.Name)
		if err != nil {
			return fmt.Errorf("ошибка создания файла %s: %v", pack.Name, err)
		}
		defer zf.Close()

		zipWriter := zip.NewWriter(zf)
		defer zipWriter.Close()

		// итерация по папкам заданным в "targets"
		for _, u := range upa.Targets {
			exclude := u.(map[string]string)["exclude"]
			// возвращает слайс имён  файлов для упаковки
			filesToZip, err := Walk(u.(map[string]string)["path"], exclude)
			if err != nil {
				return err
			}
			
			// Устанавливаем комментарий с версией
			zipWriter.SetComment(fmt.Sprintf("Version: %s", upa.Version))

			// пакуем файлы, прошедшие отбор, в архив
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
		err = ssher.LoadBySSH(models.SSHConf.Host, models.SSHConf.User, models.SSHConf.Password, pack.Name, "/"+pack.Name)
		if err != nil {
			return err
		}
	}

	return
}
