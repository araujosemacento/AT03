package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {

	var erro error
	var expressao string

	for {

		fmt.Println("Insira sua expressão:")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		expressao = scanner.Text()

		expressao, erro = higienizar(expressao)

		if erro != nil {
			fmt.Printf("\n%s\nInsira uma expressão matemática válida.\n\n", erro)
		} else {
			break
		}

	}

	fmt.Println(expressao)

}

func higienizar(expressao string) (string, error) {

	var expressaoHigienizada strings.Builder
	var simbolosPermitidos string = "()-+*/"

	for _, char := range expressao {
		if strings.Contains(simbolosPermitidos, string(char)) {

			expressaoHigienizada.WriteString(" ")
			expressaoHigienizada.WriteString(string(char))
			expressaoHigienizada.WriteString(" ")

		} else if unicode.IsSpace(char) {

			continue

		} else if unicode.IsDigit(char) || string(char) == "." {

			expressaoHigienizada.WriteString(string(char))

		} else {

			return "", fmt.Errorf("Caractere inválido na expressão: %c", char)

		}

	}

	return expressaoHigienizada.String(), nil
}
