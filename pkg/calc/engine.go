package calc

import (
	"fmt"
	"strconv"
)

// AvaliarRPN lê a fila de tokens em Notação Polonesa Reversa (RPN),
// empilha os números e chama as operações na interface [Operacoes] conforme os sinais aparecem.
func AvaliarRPN(tokens []string, ops Operacoes) (float64, error) {
	if len(tokens) == 0 {
		return 0, fmt.Errorf("expressão vazia")
	}

	var pilha []float64

	for _, token := range tokens {
		// Se o token for um número, coloca na pilha e passa para o próximo
		if num, erro := strconv.ParseFloat(token, 64); erro == nil {
			pilha = append(pilha, num)
			continue
		}

		// O sinal 'u-' inverte o sinal do último número da pilha fazendo (0 - valor)
		if token == "u-" {
			if len(pilha) < 1 {
				return 0, fmt.Errorf("expressão malformada: falta número para o sinal negativo")
			}
			topo := len(pilha) - 1
			resultado, err := ops.Subtrair(0, pilha[topo])
			if err != nil {
				return 0, err
			}
			pilha[topo] = resultado
			continue
		}

		// Operadores normais precisam de pelo menos dois números na pilha
		if len(pilha) < 2 {
			return 0, fmt.Errorf("expressão malformada: faltam números para o operador %s", token)
		}

		// Retira os dois últimos números da pilha para fazer a conta
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

		// Guarda o resultado de volta na pilha para as próximas operações
		pilha = append(pilha, resultado)
	}

	// Ao final, deve sobrar apenas o resultado final na pilha
	if len(pilha) != 1 {
		return 0, fmt.Errorf("expressão malformada: sobraram números sem operador")
	}

	return pilha[0], nil
}
