// Pacote parser disponibiliza funções para análise, higienização e conversão
// de expressões matemáticas pra Notação Polonesa Reversa (RPN).
package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Constantes com os delimitadores e operadores matemáticos aceitos pelo analisador.
const (
	Delimitadores      string = "()"
	Operadores         string = "*/+-"
	SimbolosPermitidos string = Delimitadores + Operadores
)

// Higienizar formata uma expressão matemática inserindo espaçamento entre os elementos.
//
// O objetivo dessa função é isolar operadores e parênteses com espaços laterais,
// pra manter sequências de dígitos e pontos decimais adjacentes. Preserva
// números de múltiplos dígitos e de ponto flutuante.
//
// Espaços em branco preexistentes são descartados.
//
// # Entradas
//   - expressao: texto bruto digitado pelo usuário ou recebido via rede.
//
// # Retorno
//   - string: expressão formatada com separação uniforme entre os tokens.
//   - error: retorna erro caso encontre caracteres desconhecidos (letras ou símbolos não permitidos).
func Higienizar(expressao string) (string, error) {
	var expressaoHigienizada strings.Builder

	// Itera sobre cada caractere da expressão bruta pra padronizar os espaçamentos.
	for _, char := range expressao {
		if strings.ContainsRune(SimbolosPermitidos, char) {
			// Isola operadores e parênteses com espaços pra facilitar a separação em tokens.
			expressaoHigienizada.WriteString(" ")
			expressaoHigienizada.WriteRune(char)
			expressaoHigienizada.WriteString(" ")
		} else if unicode.IsSpace(char) {
			// Descarta espaços em branco preexistentes.
			continue
		} else if unicode.IsDigit(char) || char == '.' {
			// Mantém dígitos e pontos decimais juntos pra manter tanto números inteiros quanto fracionários.
			expressaoHigienizada.WriteRune(char)
		} else {
			// Rejeita caracteres inválidos (como letras ou símbolos não suportados).
			return "", fmt.Errorf("Caractere inválido na expressão: %c", char)
		}
	}

	return expressaoHigienizada.String(), nil
}

// Tokenizar converte uma expressão higienizada em Notação Polonesa Reversa (RPN).
//
// Utiliza o algoritmo Shunting-yard pra reordenar os tokens da notação infixa
// para notação pós-fixa. A RPN elimina a necessidade de parênteses e organiza
// os operadores sequencialmente conforme sua precedência matemática.
//
// # Entradas
//   - expressaoHigienizada: texto previamente processado pela função [Higienizar].
//
// # Retorno
//   - []string: sequência linear de tokens ordenados para cálculo em pilha RPN.
//   - error: retorna erro em caso de parênteses desbalanceados ou operadores malformados.
func Tokenizar(expressaoHigienizada string) ([]string, error) {
	tokens := strings.Fields(expressaoHigienizada)
	var output []string
	var operadores []string

	// Itera sobre os tokens aplicando as regras do algoritmo Shunting-yard.
	for i, token := range tokens {
		// Se o token for um número, insere direto na fila de saída.
		if _, erro := strconv.ParseFloat(token, 64); erro == nil {
			output = append(output, token)
			continue
		}

		// Identifica se o operador é unário no início da expressão, após outro operador ou após '('
		ehUnario := i == 0 || tokens[i-1] == "(" || strings.ContainsAny(tokens[i-1], "+-*/")
		if ehUnario {
			if token == "-" {
				token = "u-"
			} else if token == "+" {
				// Unário positivo é neutro, descarta e avança
				continue
			}
		}

		// Tratando delimitadores e operadores de acordo com a precedência matemática.
		switch token {
		case "(":
			// Abre parênteses: empilha para marcar o início de uma nova subexpressão.
			operadores = append(operadores, token)
		case ")":
			// Fecha parênteses: desempilha os operadores, colocando eles na saída até encontrar o '(' correspondente.
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

			// Se a pilha esvaziar sem encontrar o '(', existem parênteses desbalanceados.
			if !encontrouPar {
				return nil, fmt.Errorf("Parênteses desbalanceados na expressão")
			}
		case "u-", "+", "-", "*", "/":
			// Desempilha operadores do topo com precedência maior ou igual antes de empilhar o operador atual.
			for len(operadores) > 0 && operadores[len(operadores)-1] != "(" && Precedencia(operadores[len(operadores)-1]) >= Precedencia(token) {
				topo := operadores[len(operadores)-1]
				operadores = operadores[:len(operadores)-1]
				output = append(output, topo)
			}
			operadores = append(operadores, token)
		default:
			// Identifica termos inválidos que não puderam ser convertidos em número nem identificados como operador.
			return nil, fmt.Errorf("Um termo inválido foi inserido na expressão: %s", token)
		}
	}

	// Descarrega os operadores que restam da pilha pra fila de saída RPN.
	for len(operadores) > 0 {
		topo := operadores[len(operadores)-1]
		operadores = operadores[:len(operadores)-1]

		// Se restar algum parêntese sem fechamento na pilha, acusa erro de sintaxe.
		if topo == "(" || topo == ")" {
			return nil, fmt.Errorf("Parênteses desbalanceados na expressão")
		}
		output = append(output, topo)
	}

	return output, nil
}

// Precedencia retorna a prioridade hierárquica de um operador aritmético.
//
// # Entradas
//   - simbolo: string representando o operador ("u-", "*", "/", "+", "-").
//
// # Retorno
//   - int: valor 3 para inversão unária; 2 para multiplicação e divisão; 1 para adição e subtração; 0 para outros símbolos.
func Precedencia(simbolo string) int {
	// Define o nível de prioridade das operações aritméticas básicas.
	switch simbolo {
	case "u-":
		// Operador unário de inversão possui a maior precedência.
		return 3
	case "*", "/":
		// Multiplicação e divisão possuem maior precedência.
		return 2
	case "+", "-":
		// Adição e subtração possuem menor precedência.
		return 1
	default:
		// Parênteses ou símbolos não avaliados como operadores diretos.
		return 0
	}
}
