package main

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

var wg sync.WaitGroup

func worker(id int, jobs <-chan string, results chan<- string, config *ssh.ClientConfig, versoes, setups []string, portas []string) {
	defer wg.Done()
	for host := range jobs {
		result := fmt.Sprintf("Worker %d processing host %s\n", id, host)
		result += processHost(host, config, versoes, setups, portas)
		results <- result
	}
}

func processHost(host string, config *ssh.ClientConfig, versoes, setups []string, portas []string) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("\nProcessando host: %s\n", host))

	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		result.WriteString(fmt.Sprintf("Erro ao conectar em %s: %v\n", host, err))
		return result.String()
	}
	defer client.Close()

	// Parar containers
	for _, versao := range versoes {
		if versao == "0" {
			break
		}
		result.WriteString(runCommand(client, fmt.Sprintf("docker stop pldpro%s", versao)))
	}

	// Verificar status
	result.WriteString(runCommand(client, "docker ps -a"))

	// Processar setups
	for qualSetup, setup := range setups {
		if setup == "0" {
			break
		}
		remotePath := path + setup
		if host != "vm-gce-000011" {
			if err := uploadFile(client, path+setup, remotePath); err != nil {
				result.WriteString(fmt.Sprintf("Erro no upload: %v\n", err))
				continue
			}
		}

		porta := portas[qualSetup]
		sedCmd := fmt.Sprintf("sed -i 's/portaInterface=..../portaInterface=%s/g' %sconfig.properties", porta, path)
		result.WriteString(runCommand(client, sedCmd))
		result.WriteString(runCommand(client, fmt.Sprintf("cd %s && sh ./%s", path, setup)))

		//cria script de inicialização
		criaArquivoBoot(client, &result)

		// Verificar status
		result.WriteString(runCommand(client, "docker ps -a"))
	}
	return result.String()
}

