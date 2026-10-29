// Cliente Socket TCP da calculadora.
// Realiza conexão TCP, envia as operações a serem feitas para o servidor e exibe a interface gráfica.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"at03/calculadora/pkg/ui"
)

// SocketClient guarda a conexão aberta e implementa [calc.Operacoes] via socket TCP síncrono.
type SocketClient struct {
	conn   net.Conn
	reader *bufio.Scanner
	writer *bufio.Writer
	mutex  sync.Mutex
}

// NovoSocketClient abre a conexão TCP com o servidor no endereço informado.
func NovoSocketClient(endereco string) (*SocketClient, error) {
	// Conecta ao servidor TCP remoto
	conn, err := net.Dial("tcp", endereco)
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao servidor em %s: %w", endereco, err)
	}

	return &SocketClient{
		conn:   conn,
		reader: bufio.NewScanner(conn),
		writer: bufio.NewWriter(conn),
	}, nil
}

// Fechar encerra a conexão com o servidor.
func (c *SocketClient) Fechar() error {
	return c.conn.Close()
}

// enviarOperacao manda a operação em texto pelo socket e espera a linha de resposta do servidor.
func (c *SocketClient) enviarOperacao(operador string, a, b float64) (float64, error) {
	// Trava o mutex para garantir que uma requisição termine antes de outra usar o mesmo socket
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Monta a linha no formato 'OP VAL1 VAL2\n'
	mensagem := fmt.Sprintf("%s %g %g\n", operador, a, b)
	if _, err := c.writer.WriteString(mensagem); err != nil {
		return 0, fmt.Errorf("erro ao escrever no socket: %w", err)
	}
	// Envia os dados do buffer para a rede na hora
	if err := c.writer.Flush(); err != nil {
		return 0, fmt.Errorf("erro ao enviar dados pelo socket: %w", err)
	}

	// Espera a linha de resposta do servidor
	if !c.reader.Scan() {
		if err := c.reader.Err(); err != nil {
			return 0, fmt.Errorf("erro ao ler resposta do socket: %w", err)
		}
		return 0, fmt.Errorf("conexão fechada pelo servidor")
	}

	resposta := strings.TrimSpace(c.reader.Text())

	// Se começar com 'OK ', converte o valor recebido para número
	if strings.HasPrefix(resposta, "OK ") {
		valorStr := strings.TrimPrefix(resposta, "OK ")
		valor, err := strconv.ParseFloat(valorStr, 64)
		if err != nil {
			return 0, fmt.Errorf("valor inválido recebido do servidor: %s", valorStr)
		}
		return valor, nil
	}

	// Se começar com 'ERRO: ', devolve a mensagem de erro vinda do servidor
	if strings.HasPrefix(resposta, "ERRO: ") {
		return 0, fmt.Errorf("%s", strings.TrimPrefix(resposta, "ERRO: "))
	}

	return 0, fmt.Errorf("resposta inesperada do servidor: %s", resposta)
}

func (c *SocketClient) Somar(a, b float64) (float64, error) {
	return c.enviarOperacao("+", a, b)
}

func (c *SocketClient) Subtrair(a, b float64) (float64, error) {
	return c.enviarOperacao("-", a, b)
}

func (c *SocketClient) Multiplicar(a, b float64) (float64, error) {
	return c.enviarOperacao("*", a, b)
}

func (c *SocketClient) Dividir(a, b float64) (float64, error) {
	return c.enviarOperacao("/", a, b)
}

func main() {
	enderecoServidor := "localhost:6001"
	log.Printf("[SOCKET] Conectando ao servidor em %s...", enderecoServidor)

	cliente, err := NovoSocketClient(enderecoServidor)
	if err != nil {
		log.Fatalf("[SOCKET] Erro de conexão: %v (rode primeiro: go run ./cmd/socket-server)", err)
	}
	defer cliente.Fechar()

	log.Println("[SOCKET] Conectado com sucesso. Abrindo a interface...")
	ui.IniciarGUI("Calculadora: Sockets TCP", cliente)
}
