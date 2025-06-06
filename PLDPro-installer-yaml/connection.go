package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func createSSHConfig() (*ssh.ClientConfig, error) {
	user, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter usuário: %w", err)
	}

	privateKeyPath := filepath.Join(user.HomeDir, ".ssh", "id_rsa")
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler chave privada: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear chave: %w", err)
	}

	return &ssh.ClientConfig{
		User: user.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

func uploadFile(client *ssh.Client, localPath, remotePath string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("erro ao criar cliente SFTP: %w", err)
	}
	defer sftpClient.Close()

	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo local: %w", err)
	}
	defer localFile.Close()

	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo remoto: %w", err)
	}
	defer remoteFile.Close()

	if _, err = io.Copy(remoteFile, localFile); err != nil {
		return fmt.Errorf("erro ao copiar arquivo: %w", err)
	}

	runCommand(client, "chmod a+x "+remotePath)
	return nil
}

func runCommand(client *ssh.Client, cmd string) string {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Sprintf("Erro ao criar sessão: %v\n", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return fmt.Sprintf("Erro no comando '%s': %v\nOutput: %s\n", cmd, err, output)
	}
	//return fmt.Sprintf("[%s] %s\nOutput:\n%s\n", client.RemoteAddr().String(), cmd, output)
	return string(output)
}

func criaArquivoBoot(client *ssh.Client, result *strings.Builder) {

	cmdString := `docker ps | grep ecomerc | grep 808 | awk '{split($2, a, ":"); print a[2]}'`
	ecomercVersion := runCommand(client, cmdString)
	result.WriteString(ecomercVersion)

	mvCmd := "mv /opt/scripts/start_container_onboot.sh /opt/scripts/start_container_onboot.sh.bkp "

	result.WriteString(runCommand(client, mvCmd))

	if len(strings.TrimSpace(ecomercVersion)) > 0 {

		shebangCmd := fmt.Sprintf(`echo "#!/bin/bash" >> /opt/scripts/start_container_onboot.sh`)
		result.WriteString(runCommand(client, shebangCmd))

		msgCmd := fmt.Sprintf(`echo "echo 'Starting PLDPro container ....'" >> /opt/scripts/start_container_onboot.sh`)
		result.WriteString(runCommand(client, msgCmd))

		for _, versao := range strings.Fields(strings.TrimSpace(ecomercVersion)) {
			criaArquivoCmd := fmt.Sprintf(`echo "cd /opt/PLDPro/ECOMERC%s && ./iniciaContainer.sh" >> /opt/scripts/start_container_onboot.sh`, versao)
			result.WriteString(runCommand(client, criaArquivoCmd))
		}
	}

	chmodCmd := "chmod a+x /opt/scripts/start_container_onboot.sh"
	result.WriteString(runCommand(client, chmodCmd))
}
