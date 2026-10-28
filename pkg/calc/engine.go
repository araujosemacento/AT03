package calc

import (
	"fmt"
	"strconv"
)

// AvaliarRPN processa uma sequência de tokens em Notação Polonesa Reversa.
//
// O algoritmo consome a lista pós-fixa utilizando uma pilha
// e delega cada operação pra interface [Operacoes].
//
// # Entradas
//   - tokens: fatia de strings contendo números e operadores ordenados em notação pós-fixa.
//   - ops: implementação do contrato [Operacoes] (instância do servidor ou cliente de rede).
//
// # Retorno
//   - float64: resultado numérico final obtido após consumir todos os tokens da expressão.
//   - error: retorna erro em caso de divisão por zero, operadores desconhecidos ou
//     expressão malformada (operandos insuficientes ou excedentes na pilha).
func AvaliarRPN(tokens []string, ops Operacoes) (float64, error) {
	if len(tokens) == 0 {
		return 0, fmt.Errorf("expressão vazia")
	}

	var pilha []float64

	// Verifica cada token da lista recebida
	for _, token := range tokens {
		// Se for um número, nenhuma excessão é acusada pela conversão de string
		// para número fluteuante, logo sendo adicionado à pilha.
		if num, erro := strconv.ParseFloat(token, 64); erro == nil {
			pilha = append(pilha, num)
			continue
		}

		// Se não houverem numerais o suficiente aguardando uma operação,
		// significa que a expressão foi enviada de forma incompleta.
		if len(pilha) < 2 {
			return 0, fmt.Errorf("Expressão malformada: operandos insuficientes para o operador %s", token)
		}

		// Retira os valores da pilha pra armazená-los em variáveis auxiliares.
		valorDois := pilha[len(pilha)-1]
		valorUm := pilha[len(pilha)-2]
		pilha = pilha[:len(pilha)-2]

		var resultado float64
		var err error

		switch token {
		case "+":
			resultado, err = ops.Somar(valorUm, valorDois)
		case "-":
			resultado, err = ops.Subtrair(valorUm, valorDois)
		case "*":
			resultado, err = ops.Multiplicar(valorUm, valorDois)
		case "/":
			resultado, err = ops.Dividir(valorUm, valorDois)
		default:
			return 0, fmt.Errorf("operador desconhecido: %s", token)
		}

		if err != nil {
			return 0, err
		}

		pilha = append(pilha, resultado)
	}

	// Após o laço, o único valor restante na pilha deve ser o resultado da expressão.
	if len(pilha) != 1 {
		return 0, fmt.Errorf("expressão malformada: quantidade incorreta de operandos")
	}

	return pilha[0], nil
}
