// Pacote calc reúne a interface e as funções que fazem os cálculos da calculadora.
package calc

import "fmt"

// Operacoes define as 4 operações básicas que qualquer cliente ou servidor precisa ter.
type Operacoes interface {
	// Somar faz a adição de dois valores (a + b).
	Somar(a, b float64) (float64, error)

	// Subtrair faz a subtração de dois valores (a - b).
	Subtrair(a, b float64) (float64, error)

	// Multiplicar faz a multiplicação de dois valores (a * b).
	Multiplicar(a, b float64) (float64, error)

	// Dividir faz a divisão de dois valores (a / b).
	Dividir(a, b float64) (float64, error)
}

// ServerOperations implementa [Operacoes] resolvendo as contas direto no servidor.
type ServerOperations struct{}

// Somar calcula a + b.
func (ServerOperations) Somar(a, b float64) (float64, error) {
	return a + b, nil
}

// Subtrair calcula a - b.
func (ServerOperations) Subtrair(a, b float64) (float64, error) {
	return a - b, nil
}

// Multiplicar calcula a * b.
func (ServerOperations) Multiplicar(a, b float64) (float64, error) {
	return a * b, nil
}

// Dividir calcula a / b e avisa se houver tentativa de divisão por zero.
func (ServerOperations) Dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("divisão por zero")
	}
	return a / b, nil
}
