// Cliente JSON-RPC da calculadora.
// Conecta na porta :6003, envia chamadas JSON-RPC para o servidor e exibe a interface gráfica.
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"

	"at03/calculadora/pkg/calc"
	"at03/calculadora/pkg/ui"
)

// RPCClient guarda o cliente conectado e implementa [calc.Operacoes] via JSON-RPC.
type RPCClient struct {
	client *rpc.Client
}

// NovoRPCClient abre a conexão TCP usando o codec JSON-RPC.
func NovoRPCClient(endereco string) (*RPCClient, error) {
	// Conecta ao servidor TCP usando JSON-RPC
	c, err := jsonrpc.Dial("tcp", endereco)
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao servidor JSON-RPC em %s: %w", endereco, err)
	}
	return &RPCClient{client: c}, nil
}

// Fechar encerra a conexão com o servidor.
func (c *RPCClient) Fechar() error {
	return c.client.Close()
}

// chamarMetodo manda a requisição em JSON e espera o resultado de forma síncrona.
func (c *RPCClient) chamarMetodo(metodo string, a, b float64) (float64, error) {
	args := calc.ArgsRPC{A: a, B: b}
	var resultado float64

	// Chama 'Calculadora.<Metodo>' serializando os argumentos em JSON
	err := c.client.Call("Calculadora."+metodo, args, &resultado)
	if err != nil {
		return 0, err
	}
	return resultado, nil
}

func (c *RPCClient) Somar(a, b float64) (float64, error) {
	return c.chamarMetodo("Somar", a, b)
}

func (c *RPCClient) Subtrair(a, b float64) (float64, error) {
	return c.chamarMetodo("Subtrair", a, b)
}

func (c *RPCClient) Multiplicar(a, b float64) (float64, error) {
	return c.chamarMetodo("Multiplicar", a, b)
}

func (c *RPCClient) Dividir(a, b float64) (float64, error) {
	return c.chamarMetodo("Dividir", a, b)
}

func main() {
	enderecoServidor := "localhost:6003"
	log.Printf("[JSON-RPC] Conectando ao servidor em %s...", enderecoServidor)

	cliente, err := NovoRPCClient(enderecoServidor)
	if err != nil {
		log.Fatalf("[JSON-RPC] Erro de conexão: %v (rode primeiro: go run ./cmd/rpc-server)", err)
	}
	defer cliente.Fechar()

	log.Println("[JSON-RPC] Conectado com sucesso. Abrindo a interface...")
	ui.IniciarGUI("Calculadora: JSON-RPC", cliente)
}
