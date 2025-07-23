package ssher

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Receiver(host, user, password, localPath, remotePath string) (err error) {
	// Параметры соединения
	//	port := "22"
	//	filePathOnServer := "/path/to/file.txt"
	//	localFilePath := "./downloaded_file.txt"

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return fmt.Errorf("не удалось подключиться: %w", err)
	}
	defer conn.Close()

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("ошибка создания SFTP-клиента: %w", err)
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть удалённый файл: %w", err)
	}
	defer remoteFile.Close()

	content, err := io.ReadAll(remoteFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	if err := os.WriteFile(localPath, content, 0644); err != nil {
		return fmt.Errorf("ошибка записи локального файла: %w", err)
	}

	fmt.Println("File downloaded successfully!")
	return
}

func DownloadBySSH(host, user, password, remotePath, localPath string) error {



	
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return fmt.Errorf("ssh dial failed: %w", err)
	}
	defer conn.Close()

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("sftp client failed: %w", err)
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.OpenFile(remotePath, os.O_RDONLY)
	if err != nil {
		return fmt.Errorf("remote file open failed: %w", err)
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("local file create failed: %w", err)
	}
	defer localFile.Close()

	// Копируем с буфером 32KB
	if _, err = io.CopyBuffer(localFile, remoteFile, make([]byte, 32*1024)); err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	// Принудительная синхронизация
	if err = localFile.Sync(); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	return nil
}
