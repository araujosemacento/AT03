// Pacote calc também inclui os tipos usados pelos servidores e clientes RPC e RMI.
package calc

import (
	"fmt"
	"log"
)

// ArgsRPC guarda os dois valores enviados em cada chamada RPC ou RMI.
type ArgsRPC struct {
	A float64
	B float64
}

// ServicoRPC adapta as operações da calculadora para o formato exigido pelo pacote net/rpc do Go.
type ServicoRPC struct {
	Ops Operacoes
}

// Somar recebe os valores via rede e preenche o ponteiro reply com o resultado de A + B.
func (s *ServicoRPC) Somar(args ArgsRPC, reply *float64) error {
	res, err := s.Ops.Somar(args.A, args.B)
	if err != nil {
		log.Printf("[RPC/RMI] Erro em Somar(%g, %g): %v", args.A, args.B, err)
		return err
	}
	log.Printf("[RPC/RMI] Operação concluída: %g + %g = %g", args.A, args.B, res)
	*reply = res
	return nil
}

// Subtrair recebe os valores via rede e preenche o ponteiro reply com o resultado de A - B.
func (s *ServicoRPC) Subtrair(args ArgsRPC, reply *float64) error {
	res, err := s.Ops.Subtrair(args.A, args.B)
	if err != nil {
		log.Printf("[RPC/RMI] Erro em Subtrair(%g, %g): %v", args.A, args.B, err)
		return err
	}
	log.Printf("[RPC/RMI] Operação concluída: %g - %g = %g", args.A, args.B, res)
	*reply = res
	return nil
}

// Multiplicar recebe os valores via rede e preenche o ponteiro reply com o resultado de A * B.
func (s *ServicoRPC) Multiplicar(args ArgsRPC, reply *float64) error {
	res, err := s.Ops.Multiplicar(args.A, args.B)
	if err != nil {
		log.Printf("[RPC/RMI] Erro em Multiplicar(%g, %g): %v", args.A, args.B, err)
		return err
	}
	log.Printf("[RPC/RMI] Operação concluída: %g * %g = %g", args.A, args.B, res)
	*reply = res
	return nil
}

// Dividir recebe os valores via rede e preenche o ponteiro reply com o resultado de A / B.
func (s *ServicoRPC) Dividir(args ArgsRPC, reply *float64) error {
	res, err := s.Ops.Dividir(args.A, args.B)
	if err != nil {
		log.Printf("[RPC/RMI] Erro em Dividir(%g, %g): %v", args.A, args.B, err)
		return fmt.Errorf("divisão por zero")
	}
	log.Printf("[RPC/RMI] Operação concluída: %g / %g = %g", args.A, args.B, res)
	*reply = res
	return nil
}
