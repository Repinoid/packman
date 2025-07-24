package functions

import (
	"archive/zip"
	"errors"
	"fmt"
	"gorcom/internal/models"
	"gorcom/internal/ssher"
	"log"
	"os"
)

func UnPack(upa models.Packages) (err error) {

	// если у какого то пакета нет имени - на выход
	for _, packa := range upa.Packs {
		if packa.Name == "" {
			return errors.New("no \"Name\" field on Packets")
		}
	}
	for _, packa := range upa.Packs {

		recFolder := "fromSSH"
		if _, err := os.Stat(recFolder); os.IsNotExist(err) {
			err := os.Mkdir(recFolder, 0777)
			if err != nil {
				return err
			}
		}
		
		tmpFile := packa.Name + ".tmp"
		distantFile := "/files/" + packa.Name
		receivedFile := recFolder + "/New_" + packa.Name

		err = ssher.Receiver(models.SSHConf.Host, models.SSHConf.User, models.SSHConf.Password, tmpFile, distantFile)
		if err != nil {
			return err
		}
		// если в пакете задано условие по версии
		if packa.Ver != "" {
			zipReader, err := zip.OpenReader(tmpFile)
			if err != nil {
				log.Fatal(err)
			}
			// Получаем комментарий, в котором прописана версия полученного пакета
			comment := zipReader.Comment
			zipReader.Close()
			op, right, err := ParseComparisonWithRegex(packa.Ver)
			if err != nil {
				return fmt.Errorf("ошибка условия версии  %s: %v", comment, err)
			}
			// если не выполняется условие, заданное в версии пакета, удаляем скаченный файл
			if !compara(comment, op, right) {
				os.Remove(tmpFile)
			} else {
				os.Rename(tmpFile, receivedFile)
			}
		} else {
			err = os.Rename(tmpFile, receivedFile)
			if err != nil {
				return fmt.Errorf("oшибка при переименовании %s to %s, err %w", tmpFile, receivedFile, err)
			}
		}
	}
	return
}
