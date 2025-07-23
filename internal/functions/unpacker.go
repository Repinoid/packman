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
	for _, u := range upa.Packs {
		if u.Name == "" {
			return errors.New("no \"Name\" field on Packets")
		}
	}
	for _, u := range upa.Packs {
		err = ssher.Receiver(models.SSHConf.Host, models.SSHConf.User, models.SSHConf.Password, u.Name+".tmp", "/"+u.Name)
		if err != nil {
			return err
		}
		// если в пакете задано условие по версии
		if u.Ver != "" {
			zipReader, err := zip.OpenReader(u.Name + ".tmp")
			if err != nil {
				log.Fatal(err)
			}
			// Получаем комментарий, в котором прописана версия полученного пакета
			comment := zipReader.Comment
			zipReader.Close()
			op, right, err := ParseComparisonWithRegex(u.Ver)
			if err != nil {
				return fmt.Errorf("ошибка условия версии  %s: %v", comment, err)
			}
			// если не выполняется условие, заданное в версии пакета, удаляем скаченный файл
			if !compara(comment, op, right) {
				os.Remove(u.Name + ".tmp")
			} else {
				os.Rename(u.Name+".tmp", "New_"+u.Name)
			}
		} else {
			err = os.Rename(u.Name+".tmp", "New_"+u.Name)
			if err != nil {
				fmt.Println("Ошибка при переименовании:", err)
			}
		}
	}
	return
}
