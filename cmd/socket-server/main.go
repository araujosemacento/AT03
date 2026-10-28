// Servidor Socket TCP para a calculadora.
//
// Este binário implementa o servidor remoto que escuta conexões TCP na porta :6001,
// processa requisições de operações aritméticas básicas enviadas pelos clientes
// de forma concorrente e devolve as respostas formatadas em linha de texto.
package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"at03/calculadora/pkg/calc"
)

func main() {
	// Define a porta TCP padrão para escutar conexões de clientes
	endereco := ":6001"

	// Cria o listener TCP no endereço especificado com a biblioteca net nativa
	listener, err := net.Listen("tcp", endereco)
	if err != nil {
		fmt.Printf("Erro ao iniciar o servidor socket: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Servidor Socket TCP da Calculadora ativo na porta %s\n", endereco)

	// Instancia o executor de operações aritméticas do servidor
	ops := calc.ServerOperations{}

	// Laço contínuo para aceitar novas conexões de clientes
	for {
		// Bloqueia a execução até que um novo cliente solicite conexão
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Erro ao aceitar conexão: %v\n", err)
			continue
		}

		// Inicializa uma goroutine concorrente para atender o cliente de forma isolada
		go tratarCliente(conn, ops)
	}
}

// tratarCliente gerencia a comunicação bidirecional com um cliente conectado via socket TCP.
//
// # Entradas
//   - conn: conexão de rede [net.Conn] estabelecida com o cliente.
//   - ops: provedor de cálculo [calc.Operacoes] utilizado para executar as operações.
//
// # Protocolo de Mensagens
//   - Requisição do cliente: "OP VAL1 VAL2\n" (exemplo: "+ 10 20\n")
//   - Resposta com sucesso: "OK <resultado>\n" (exemplo: "OK 30\n")
//   - Resposta com falha: "ERRO: <mensagem>\n" (exemplo: "ERRO: divisão por zero\n")
func tratarCliente(conn net.Conn, ops calc.Operacoes) {
	// Garante o fechamento do socket ao encerrar o atendimento do cliente
	defer conn.Close()

	fmt.Printf("Cliente conectado: %s\n", conn.RemoteAddr().String())

	// Inicializa o leitor com buffer para receber linhas da conexão
	scanner := bufio.NewScanner(conn)
	// Inicializa o escritor com buffer para enviar respostas pela conexão
	writer := bufio.NewWriter(conn)

	// Lê linha por linha transmitida pelo cliente via TCP
	for scanner.Scan() {
		linha := strings.TrimSpace(scanner.Text())
		if linha == "" {
			continue
		}

		// Divide a mensagem em operador e dois operandos
		partes := strings.Fields(linha)
		if len(partes) != 3 {
			writer.WriteString("ERRO: formato esperado 'OP VAL1 VAL2'\n")
			writer.Flush()
			continue
		}

		operador := partes[0]
		a, errA := strconv.ParseFloat(partes[1], 64)
		b, errB := strconv.ParseFloat(partes[2], 64)

		if errA != nil || errB != nil {
			writer.WriteString("ERRO: operandos numéricos inválidos\n")
			writer.Flush()
			continue
		}

		var resultado float64
		var opErr error

		// Executa a operação aritmética solicitada
		switch operador {
		case "+":
			resultado, opErr = ops.Somar(a, b)
		case "-":
			resultado, opErr = ops.Subtrair(a, b)
		case "*":
			resultado, opErr = ops.Multiplicar(a, b)
		case "/":
			resultado, opErr = ops.Dividir(a, b)
		default:
			opErr = fmt.Errorf("operador desconhecido: %s", operador)
		}

		// Formata e envia a resposta correspondente pelo socket
		if opErr != nil {
			writer.WriteString(fmt.Sprintf("ERRO: %s\n", opErr.Error()))
		} else {
			writer.WriteString(fmt.Sprintf("OK %g\n", resultado))
		}
		// Descarrega o buffer para assegurar a transmissão imediata dos bytes pela rede
		writer.Flush()
	}

	// Verifica se a saída do laço ocorreu por falha de leitura ou erro de rede
	if err := scanner.Err(); err != nil {
		fmt.Printf("Erro na leitura do socket (%s): %v\n", conn.RemoteAddr().String(), err)
	}

	fmt.Printf("Conexão encerrada: %s\n", conn.RemoteAddr().String())
}
