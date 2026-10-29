// Pacote parser reúne funções para limpar, validar e converter
// expressões matemáticas para Notação Polonesa Reversa (RPN).
package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	Delimitadores      string = "()"
	Operadores         string = "*/+-"
	SimbolosPermitidos string = Delimitadores + Operadores
)

// Higienizar adiciona espaços ao redor de operadores e parênteses,
// mantém números inteiros ou decimais juntos e remove espaços extras.
func Higienizar(expressao string) (string, error) {
	var expressaoHigienizada strings.Builder

	for _, char := range expressao {
		if strings.ContainsRune(SimbolosPermitidos, char) {
			// Isola operadores e parênteses com espaços para facilitar a separação
			expressaoHigienizada.WriteString(" ")
			expressaoHigienizada.WriteRune(char)
			expressaoHigienizada.WriteString(" ")
		} else if unicode.IsSpace(char) {
			// Ignora espaços que já existiam na expressão
			continue
		} else if unicode.IsDigit(char) || char == '.' {
			// Mantém dígitos e pontos decimais unidos no mesmo número
			expressaoHigienizada.WriteRune(char)
		} else {
			return "", fmt.Errorf("caractere inválido na expressão: %c", char)
		}
	}

	return expressaoHigienizada.String(), nil
}

// Tokenizar transforma a expressão higienizada em uma fila RPN (Notação Polonesa Reversa)
// usando o algoritmo Shunting-yard para organizar a ordem das operações e parênteses.
func Tokenizar(expressaoHigienizada string) ([]string, error) {
	tokens := strings.Fields(expressaoHigienizada)
	var output []string
	var operadores []string

	for i, token := range tokens {
		// Se for número, vai direto para a fila de saída
		if _, erro := strconv.ParseFloat(token, 64); erro == nil {
			output = append(output, token)
			continue
		}

		// Verifica se o sinal é unário (no começo, após '(' ou após outro operador)
		ehUnario := i == 0 || tokens[i-1] == "(" || strings.ContainsAny(tokens[i-1], "+-*/")
		if ehUnario {
			if token == "-" {
				token = "u-"
			} else if token == "+" {
				// Sinal de '+' unário não muda o valor, então é descartado
				continue
			}
		}

		switch token {
		case "(":
			operadores = append(operadores, token)
		case ")":
			// Desempilha tudo até achar o '(' que abriu o bloco
			encontrouPar := false
			for len(operadores) > 0 {
				topo := operadores[len(operadores)-1]
				operadores = operadores[:len(operadores)-1]

				if topo == "(" {
					encontrouPar = true
					break
				}
				output = append(output, topo)
			}

			if !encontrouPar {
				return nil, fmt.Errorf("parênteses desbalanceados na expressão")
			}
		case "u-", "+", "-", "*", "/":
			// O operador unário 'u-' agrupa da direita para a esquerda (desempilha só se topo > token)
			// Os operadores normais (+, -, *, /) agrupam da esquerda para a direita (desempilha se topo >= token)
			for len(operadores) > 0 && operadores[len(operadores)-1] != "(" {
				topo := operadores[len(operadores)-1]
				deveDesempilhar := false
				if token == "u-" {
					deveDesempilhar = Precedencia(topo) > Precedencia(token)
				} else {
					deveDesempilhar = Precedencia(topo) >= Precedencia(token)
				}

				if !deveDesempilhar {
					break
				}

				operadores = operadores[:len(operadores)-1]
				output = append(output, topo)
			}
			operadores = append(operadores, token)
		default:
			return nil, fmt.Errorf("termo inválido na expressão: %s", token)
		}
	}

	// Passa o restante dos operadores da pilha para a fila de saída
	for len(operadores) > 0 {
		topo := operadores[len(operadores)-1]
		operadores = operadores[:len(operadores)-1]

		if topo == "(" || topo == ")" {
			return nil, fmt.Errorf("parênteses desbalanceados na expressão")
		}
		output = append(output, topo)
	}

	return output, nil
}

// Precedencia define qual operador deve ser resolvido primeiro (quanto maior o valor, maior a prioridade).
func Precedencia(simbolo string) int {
	switch simbolo {
	case "u-":
		return 3
	case "*", "/":
		return 2
	case "+", "-":
		return 1
	default:
		return 0
	}
}
