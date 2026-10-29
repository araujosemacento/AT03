// Servidor RMI (Remote Method Invocation) em Go usando net/rpc e Gob.
// Fica escutando na porta :6002 e deixa o objeto Calculadora disponível para chamadas remotas.
package main

import (
	"log"
	"net"
	"net/rpc"

	"at03/calculadora/pkg/calc"
)

func main() {
	endereco := ":6002"

	// Cria e registra o objeto Calculadora no servidor RPC do Go
	servico := &calc.ServicoRPC{Ops: calc.ServerOperations{}}
	if err := rpc.RegisterName("Calculadora", servico); err != nil {
		log.Fatalf("[RMI] Erro ao registrar o objeto Calculadora: %v", err)
	}

	// Abre a porta TCP para receber as conexões dos clientes
	listener, err := net.Listen("tcp", endereco)
	if err != nil {
		log.Fatalf("[RMI] Erro ao iniciar o servidor na porta %s: %v", endereco, err)
	}
	defer listener.Close()

	log.Printf("[RMI] Servidor rodando e esperando chamadas na porta %s", endereco)

	// Fica em loop aceitando clientes e atendendo cada um numa goroutine nova
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[RMI] Erro ao aceitar conexão: %v", err)
			continue
		}

		log.Printf("[RMI] Cliente conectado: %s", conn.RemoteAddr().String())
		// Atende as chamadas remotas do cliente usando serialização binária Gob
		go rpc.ServeConn(conn)
	}
}
