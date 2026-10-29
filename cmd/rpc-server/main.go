// Servidor JSON-RPC da calculadora.
// Fica escutando na porta :6003 e responde às requisições em formato JSON sobre TCP.
package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"at03/calculadora/pkg/calc"
)

func main() {
	endereco := ":6003"

	// Cria um servidor RPC e registra os métodos de Calculadora
	servico := &calc.ServicoRPC{Ops: calc.ServerOperations{}}
	server := rpc.NewServer()
	if err := server.RegisterName("Calculadora", servico); err != nil {
		log.Fatalf("[JSON-RPC] Erro ao registrar o serviço: %v", err)
	}

	// Abre a porta TCP para esperar conexões
	listener, err := net.Listen("tcp", endereco)
	if err != nil {
		log.Fatalf("[JSON-RPC] Erro ao abrir a porta %s: %v", endereco, err)
	}
	defer listener.Close()

	log.Printf("[JSON-RPC] Servidor rodando e esperando conexões na porta %s", endereco)

	// Fica em loop aceitando clientes e atendendo cada um numa goroutine nova
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[JSON-RPC] Erro ao aceitar conexão: %v", err)
			continue
		}

		log.Printf("[JSON-RPC] Cliente conectado: %s", conn.RemoteAddr().String())
		// Lê e responde às chamadas usando o codec JSON-RPC
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
