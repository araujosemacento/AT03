// Cliente RMI da calculadora.
// Conecta via net/rpc na porta :6002, chama os métodos remotos do servidor e abre a interface gráfica.
package main

import (
	"fmt"
	"log"
	"net/rpc"

	"at03/calculadora/pkg/calc"
	"at03/calculadora/pkg/ui"
)

// RMIClient guarda o cliente RPC conectado e implementa [calc.Operacoes].
type RMIClient struct {
	client *rpc.Client
}

// NovoRMIClient conecta ao servidor RMI pelo endereço TCP informado.
func NovoRMIClient(endereco string) (*RMIClient, error) {
	// Abre a conexão RPC padrão do Go (usa Gob sobre TCP)
	c, err := rpc.Dial("tcp", endereco)
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao servidor RMI em %s: %w", endereco, err)
	}
	return &RMIClient{client: c}, nil
}

// Fechar encerra a conexão com o servidor RMI.
func (c *RMIClient) Fechar() error {
	return c.client.Close()
}

// chamarMetodo invoca um método do objeto remoto Calculadora de forma síncrona.
func (c *RMIClient) chamarMetodo(metodo string, a, b float64) (float64, error) {
	args := calc.ArgsRPC{A: a, B: b}
	var resultado float64

	// Chama 'Calculadora.<Metodo>' no servidor passando os valores e esperando o resultado
	err := c.client.Call("Calculadora."+metodo, args, &resultado)
	if err != nil {
		return 0, err
	}
	return resultado, nil
}

func (c *RMIClient) Somar(a, b float64) (float64, error) {
	return c.chamarMetodo("Somar", a, b)
}

func (c *RMIClient) Subtrair(a, b float64) (float64, error) {
	return c.chamarMetodo("Subtrair", a, b)
}

func (c *RMIClient) Multiplicar(a, b float64) (float64, error) {
	return c.chamarMetodo("Multiplicar", a, b)
}

func (c *RMIClient) Dividir(a, b float64) (float64, error) {
	return c.chamarMetodo("Dividir", a, b)
}

func main() {
	enderecoServidor := "localhost:6002"
	log.Printf("[RMI] Conectando ao servidor em %s...", enderecoServidor)

	cliente, err := NovoRMIClient(enderecoServidor)
	if err != nil {
		log.Fatalf("[RMI] Erro de conexão: %v (rode primeiro: go run ./cmd/rmi-server)", err)
	}
	defer cliente.Fechar()

	log.Println("[RMI] Conectado com sucesso. Abrindo a interface...")
	ui.IniciarGUI("Calculadora: RMI (Go RPC)", cliente)
}
