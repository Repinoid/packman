package functions

import (
	"errors"
	"gorcom/internal/models"
	"gorcom/internal/ssher"
)

func UnPack(upa models.Packages) (err error) {

	// если у какого то пакета нет имени - на выход
	for _, u := range upa.Packs {
		if u.Name == "" {
			return errors.New("no \"Name\" field on Packets")
		}
	}
	for _, u := range upa.Packs {
		err = ssher.Receiver(models.SSHConf.Host, models.SSHConf.User, models.SSHConf.Password, u.Name, "/" + u.Name)
		if err != nil {
			return err
		}

		if u.Ver == "" {

		}
	}

	//receiver(host, user, password, localPath, remotePath string) {
	return
}
