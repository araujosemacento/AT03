// Servidor Socket TCP da calculadora.
// Fica escutando na porta :6001, recebe as operações enviadas pelos clientes e devolve os resultados.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"at03/calculadora/pkg/calc"
)

func main() {
	endereco := ":6001"

	// Abre o socket TCP para esperar conexões na porta definida
	listener, err := net.Listen("tcp", endereco)
	if err != nil {
		log.Fatalf("[SOCKET] Erro ao iniciar o servidor na porta %s: %v", endereco, err)
	}
	defer listener.Close()

	log.Printf("[SOCKET] Servidor rodando e esperando conexões na porta %s", endereco)

	ops := calc.ServerOperations{}

	// Fica em loop aceitando novos clientes que se conectarem
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[SOCKET] Erro ao aceitar conexão: %v", err)
			continue
		}

		// Atende cada cliente ao mesmo tempo numa goroutine nova pra cada um
		go tratarCliente(conn, ops)
	}
}

// tratarCliente lê as linhas no formato 'OP VAL1 VAL2\n', faz a conta e responde pelo socket.
func tratarCliente(conn net.Conn, ops calc.Operacoes) {
	defer conn.Close()

	clienteRemoto := conn.RemoteAddr().String()
	log.Printf("[SOCKET] Cliente conectado: %s", clienteRemoto)

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	// Lê cada linha que chega do cliente enquanto a conexão estiver aberta
	for scanner.Scan() {
		linha := strings.TrimSpace(scanner.Text())
		if linha == "" {
			continue
		}

		// Quebra a linha em 3 partes: operador, primeiro valor e segundo valor
		partes := strings.Fields(linha)
		if len(partes) != 3 {
			log.Printf("[SOCKET] [%s] Mensagem fora do formato esperado: %s", clienteRemoto, linha)
			writer.WriteString("ERRO: formato esperado 'OP VAL1 VAL2'\n")
			writer.Flush()
			continue
		}

		operador := partes[0]
		a, errA := strconv.ParseFloat(partes[1], 64)
		b, errB := strconv.ParseFloat(partes[2], 64)

		if errA != nil || errB != nil {
			log.Printf("[SOCKET] [%s] Valores inválidos recebidos: %s", clienteRemoto, linha)
			writer.WriteString("ERRO: valores numéricos inválidos\n")
			writer.Flush()
			continue
		}

		var resultado float64
		var opErr error

		log.Printf("[SOCKET] [%s] Operação pedida: %g %s %g", clienteRemoto, a, operador, b)

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

		// Escreve o resultado ou o erro de volta para o cliente
		if opErr != nil {
			log.Printf("[SOCKET] [%s] Erro na operação: %v", clienteRemoto, opErr)
			writer.WriteString(fmt.Sprintf("ERRO: %s\n", opErr.Error()))
		} else {
			log.Printf("[SOCKET] [%s] Resultado enviado: %g %s %g = %g", clienteRemoto, a, operador, b, resultado)
			writer.WriteString(fmt.Sprintf("OK %g\n", resultado))
		}
		// Envia os dados do buffer para a rede na hora
		writer.Flush()
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[SOCKET] [%s] Erro ao ler dados do socket: %v", clienteRemoto, err)
	}

	log.Printf("[SOCKET] Conexão fechada com %s", clienteRemoto)
}
