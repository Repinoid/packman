package ssher

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func createSSHClient(host string, user string, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Это небезопасная настройка! Используйте правильный callback!
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, fmt.Errorf("не удалось установить соединение: %w", err)
	}

	return client, nil
}

func uploadFile(sftpClient *sftp.Client, localPath string, remotePath string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия локального файла %s: %w", localPath, err)
	}
	defer localFile.Close()

	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("ошибка создания удалённого файла %s: %w", remotePath, err)
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("ошибка копирования файла %s to %s: %v", remoteFile.Name(), localFile.Name(), err)
	}

	fmt.Printf("Файл %s успешно загружен на сервер как %s\n", localPath, remotePath)
	return nil
}

func LoadBySSH(host, user, password, localPath, remotePath string) (err error) {

	// Подключаемся по SSH
	sshClient, err := createSSHClient(host, user, password)
	if err != nil {
		return fmt.Errorf("createSSHClient: %w", err)
	}
	defer sshClient.Close()

	// Получаем SFTP-клиент
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("oшибка инициализации SFTP клиента: %w", err)
	}
	defer sftpClient.Close()

	// Загружаем файл на сервер
	err = uploadFile(sftpClient, localPath, remotePath)
	return
}
