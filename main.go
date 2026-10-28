package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const Delimitadores string = "()"
const Operadores string = "*/+-"
const SimbolosPermitidos string = Delimitadores + Operadores

func main() {

	var erro error
	var expressao string
	var tokens []string

	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Println("Insira sua expressão:")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Println("Erro na leitura da entrada:", err)
			}
			return
		}
		expressao = scanner.Text()

		expressao, erro = higienizar(expressao)
		if erro != nil {
			fmt.Printf("\n%s\nInsira uma expressão matemática válida.\n\n", erro)
			continue
		}

		tokens, erro = tokenizar(expressao)
		if erro != nil {
			fmt.Printf("\n%s\nInsira uma expressão matemática válida.\n\n", erro)
			continue
		}

		break

	}

	fmt.Println(tokens)

}

func higienizar(expressao string) (string, error) {

	var expressaoHigienizada strings.Builder

	for _, char := range expressao {
		if strings.ContainsRune(SimbolosPermitidos, char) {

			expressaoHigienizada.WriteString(" ")
			expressaoHigienizada.WriteRune(char)
			expressaoHigienizada.WriteString(" ")

		} else if unicode.IsSpace(char) {

			continue

		} else if unicode.IsDigit(char) || char == '.' {

			expressaoHigienizada.WriteRune(char)

		} else {

			return "", fmt.Errorf("Caractere inválido na expressão: %c", char)

		}

	}

	return expressaoHigienizada.String(), nil
}

func tokenizar(expressaoHigienizada string) ([]string, error) {

	tokens := strings.Fields(expressaoHigienizada)
	var output []string
	var operadores []string

	for _, token := range tokens {
		if _, erro := strconv.ParseFloat(token, 64); erro == nil {
			output = append(output, token)
			continue
		}

		switch token {
		case "(":
			operadores = append(operadores, token)
		case ")":
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
				return nil, fmt.Errorf("Parênteses desbalanceados na expressão")
			}
		case "+", "-", "*", "/":
			for len(operadores) > 0 && operadores[len(operadores)-1] != "(" && precedencia(operadores[len(operadores)-1]) >= precedencia(token) {
				topo := operadores[len(operadores)-1]
				operadores = operadores[:len(operadores)-1]
				output = append(output, topo)
			}
			operadores = append(operadores, token)
		default:
			return nil, fmt.Errorf("Um termo inválido foi inserido na expressão: %s", token)
		}
	}

	for len(operadores) > 0 {
		topo := operadores[len(operadores)-1]
		operadores = operadores[:len(operadores)-1]

		if topo == "(" || topo == ")" {
			return nil, fmt.Errorf("Parênteses desbalanceados na expressão")
		}
		output = append(output, topo)
	}

	return output, nil

}

func precedencia(simbolo string) int {
	switch simbolo {
	case "*", "/":
		return 2
	case "+", "-":
		return 1
	default:
		return 0
	}
}

// func calcular() {}
