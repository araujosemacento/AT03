// Cliente Socket TCP para a calculadora.
//
// Este binário implementa o cliente remoto que estabelece uma conexão TCP com
// o servidor socket na porta :6001, implementa o contrato [calc.Operacoes]
// transmitindo mensagens de rede e inicializa a interface gráfica.
package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"at03/calculadora/pkg/ui"
)

// SocketClient gerencia a conexão de rede e implementa a interface [calc.Operacoes] via TCP.
type SocketClient struct {
	conn   net.Conn
	reader *bufio.Scanner
	writer *bufio.Writer
	mutex  sync.Mutex
}

// NovoSocketClient cria e conecta um novo cliente socket TCP ao endereço.
//
// # Entradas
//   - endereco: host e porta do servidor TCP.
//
// # Retorno
//   - *SocketClient: ponteiro para o cliente conectado e inicializado.
//   - error: retorna erro caso a conexão TCP não possa ser estabelecida.
func NovoSocketClient(endereco string) (*SocketClient, error) {
	// Abre a conexão TCP com o servidor remoto
	conn, err := net.Dial("tcp", endereco)
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao servidor socket em %s: %w", endereco, err)
	}

	return &SocketClient{
		conn:   conn,
		reader: bufio.NewScanner(conn),
		writer: bufio.NewWriter(conn),
	}, nil
}

// Fechar encerra a conexão de rede ativa com o servidor.
func (c *SocketClient) Fechar() error {
	return c.conn.Close()
}

// enviarOperacao transmite uma operação aritmética ao servidor e aguarda a resposta síncrona.
//
// Utiliza um mutex para garantir sincronia entre envio da requisição e leitura
// da resposta através do mesmo socket.
//
// # Entradas
//   - operador: símbolo aritmético da operação ("+", "-", "*", "/").
//   - a: primeiro operando numérico.
//   - b: segundo operando numérico.
//
// # Retorno
//   - float64: resultado numérico calculado pelo servidor remoto.
//   - error: retorna erro em caso de falha de conexão, timeout ou erro aritmético retornado pelo servidor.
func (c *SocketClient) enviarOperacao(operador string, a, b float64) (float64, error) {
	// Bloqueia o acesso concorrente ao socket para sincronizar requisição e resposta
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Formata a mensagem conforme o protocolo: 'OP VAL1 VAL2\n'
	mensagem := fmt.Sprintf("%s %g %g\n", operador, a, b)
	if _, err := c.writer.WriteString(mensagem); err != nil {
		return 0, fmt.Errorf("falha ao enviar dados pelo socket: %w", err)
	}
	// Força o envio imediato dos bytes através da rede
	if err := c.writer.Flush(); err != nil {
		return 0, fmt.Errorf("falha ao descarregar buffer do socket: %w", err)
	}

	// Aguarda a resposta em linha do servidor pelo socket
	if !c.reader.Scan() {
		if err := c.reader.Err(); err != nil {
			return 0, fmt.Errorf("falha ao ler resposta do socket: %w", err)
		}
		return 0, fmt.Errorf("conexão encerrada prematuramente pelo servidor")
	}

	resposta := strings.TrimSpace(c.reader.Text())

	// Trata resposta de sucesso ('OK <resultado>')
	if strings.HasPrefix(resposta, "OK ") {
		valorStr := strings.TrimPrefix(resposta, "OK ")
		valor, err := strconv.ParseFloat(valorStr, 64)
		if err != nil {
			return 0, fmt.Errorf("resposta numérica inválida do servidor: %s", valorStr)
		}
		return valor, nil
	}

	// Trata resposta de erro ('ERRO: <mensagem>')
	if strings.HasPrefix(resposta, "ERRO: ") {
		return 0, fmt.Errorf("%s", strings.TrimPrefix(resposta, "ERRO: "))
	}

	return 0, fmt.Errorf("resposta inesperada do servidor: %s", resposta)
}

// Somar envia uma requisição de adição remota ao servidor socket.
func (c *SocketClient) Somar(a, b float64) (float64, error) {
	return c.enviarOperacao("+", a, b)
}

// Subtrair envia uma requisição de subtração remota ao servidor socket.
func (c *SocketClient) Subtrair(a, b float64) (float64, error) {
	return c.enviarOperacao("-", a, b)
}

// Multiplicar envia uma requisição de multiplicação remota ao servidor socket.
func (c *SocketClient) Multiplicar(a, b float64) (float64, error) {
	return c.enviarOperacao("*", a, b)
}

// Dividir envia uma requisição de divisão remota ao servidor socket.
func (c *SocketClient) Dividir(a, b float64) (float64, error) {
	return c.enviarOperacao("/", a, b)
}

func main() {
	enderecoServidor := "localhost:6001"
	fmt.Printf("Conectando ao servidor Socket TCP em %s...\n", enderecoServidor)

	// Estabelece a conexão com o servidor de sockets
	cliente, err := NovoSocketClient(enderecoServidor)
	if err != nil {
		fmt.Printf("Erro de conexão: %v\n", err)
		fmt.Println("Inicie primeiro o servidor executando: go run ./cmd/socket-server")
		return
	}
	defer cliente.Fechar()

	fmt.Println("Conexão estabelecida com sucesso. Abrindo interface gráfica...")

	// Inicia a interface gráfica Fyne passando o cliente socket como provedor das operações
	ui.IniciarGUI("Calculadora: Sockets TCP", cliente)
}
