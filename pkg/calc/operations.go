// Pacote calc disponibiliza contratos e motores de cálculo para as expressões
// matemáticas, suportando a execução na instância do servidor.
package calc

import "fmt"

// Operacoes define o contrato de cálculo aritmético que os clientes distribuídos
// e as instâncias de servidor devem atender.
//
// O objetivo dessa interface é desacoplar o motor de avaliação de expressões e a interface gráfica
// dos detalhes de rede ou do protocolo usado (Sockets TCP, HTTP-REST, RPC, etc.).
type Operacoes interface {
	// Somar executa a adição de dois valores (a + b).
	//
	// Pra então devolver o resultado através da rede. 🤡
	//
	// Acabou de cair a ficha do quão desnecessário isso é.
	//
	// # Entradas
	//   - a: primeiro termo numérico.
	//   - b: segundo termo numérico.
	//
	// # Retorno
	//   - float64: resultado da soma.
	//   - error: retorna erro caso ocorra falha de rede ou de processamento no servidor.
	Somar(a, b float64) (float64, error)

	// Subtrair executa a subtração entre dois valores (a - b).
	//
	// # Entradas
	//   - a: minuendo.
	//   - b: subtraendo.
	//
	// # Retorno
	//   - float64: resultado da diferença.
	//   - error: retorna erro caso ocorra falha de rede ou de processamento no servidor.
	Subtrair(a, b float64) (float64, error)

	// Multiplicar executa o produto aritmético entre dois valores (a * b).
	//
	// # Entradas
	//   - a: primeiro fator.
	//   - b: segundo fator.
	//
	// # Retorno
	//   - float64: resultado da multiplicação.
	//   - error: retorna erro caso ocorra falha de rede ou de processamento no servidor.
	Multiplicar(a, b float64) (float64, error)

	// Dividir executa o quociente aritmético entre dois valores (a / b).
	//
	// # Entradas
	//   - a: dividendo.
	//   - b: divisor.
	//
	// # Retorno
	//   - float64: quociente da divisão.
	//   - error: retorna erro caso o divisor seja zero ou ocorra falha de rede no servidor.
	Dividir(a, b float64) (float64, error)
}

// ServerOperations implementa a interface [Operacoes] executando os cálculos na instância do servidor.
type ServerOperations struct{}

// Somar calcula a soma entre a e b.
//
// # Entradas
//   - a: primeiro termo numérico.
//   - b: segundo termo numérico.
//
// # Retorno
//   - float64: resultado da soma a + b.
//   - error: sempre nil.
func (ServerOperations) Somar(a, b float64) (float64, error) {
	return a + b, nil
}

// Subtrair calcula a subtração entre a e b.
//
// # Entradas
//   - a: minuendo.
//   - b: subtraendo.
//
// # Retorno
//   - float64: resultado da diferença a - b.
//   - error: sempre nil.
func (ServerOperations) Subtrair(a, b float64) (float64, error) {
	return a - b, nil
}

// Multiplicar calcula o produto entre a e b.
//
// # Entradas
//   - a: primeiro fator.
//   - b: segundo fator.
//
// # Retorno
//   - float64: resultado do produto a * b.
//   - error: sempre nil.
func (ServerOperations) Multiplicar(a, b float64) (float64, error) {
	return a * b, nil
}

// Dividir calcula a divisão de a por b.
//
// # Entradas
//   - a: dividendo.
//   - b: divisor.
//
// # Retorno
//   - float64: quociente da divisão a / b.
//   - error: retorna erro caso o divisor b seja igual a zero.
func (ServerOperations) Dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("divisão por zero")
	}
	return a / b, nil
}
