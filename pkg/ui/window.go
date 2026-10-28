// Package ui disponibiliza componentes gráficos para a calculadora.
package ui

import (
	"fmt"
	"strings"

	"at03/calculadora/pkg/calc"
	"at03/calculadora/pkg/parser"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// IniciarGUI inicializa e exibe a interface gráfica da calculadora conectada ao provedor de operações.
//
// O objetivo dessa função é instanciar a janela usando a biblioteca Fyne,
// configurar o mostrador (histórico superior e campo de entrada inferior),
// conectar os eventos de teclado físico e pinpad, e executar a decomposição
// das expressões com envio de operações ao contrato [calc.Operacoes].
//
// # Entradas
//   - titulo: texto exibido na barra superior da janela gráfica.
//   - ops: implementação de [calc.Operacoes].
//
// # Comportamento
//
// Esta chamada é bloqueante, assumindo o laço de eventos da interface gráfica
// até que a janela seja fechada pelo usuário.
func IniciarGUI(titulo string, ops calc.Operacoes) {
	// Criação da instância da janela principal
	aplicacao := app.New()
	janela := aplicacao.NewWindow(titulo)
	janela.Resize(fyne.NewSize(320, 420))

	// Rótulo superior para exibir o histórico da expressão ou mensagens de erro
	historicoLabel := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})

	// Campo de texto principal pra entrada de dados
	displayEntry := widget.NewEntry()
	displayEntry.SetPlaceHolder("0")

	limparNoProximoInput := false

	inserirTexto := func(s string) {
		if limparNoProximoInput {
			displayEntry.SetText("")
			limparNoProximoInput = false
		}
		displayEntry.SetText(displayEntry.Text + s)
	}

	// Reseta todo o estado da interface (botão C)
	limparTudo := func() {
		historicoLabel.SetText("")
		displayEntry.SetText("")
		limparNoProximoInput = false
	}

	// Remove o último caractere digitado (botão DEL / Backspace)
	apagarUltimo := func() {
		if limparNoProximoInput {
			displayEntry.SetText("")
			limparNoProximoInput = false
			return
		}
		texto := displayEntry.Text
		if len(texto) > 0 {
			displayEntry.SetText(texto[:len(texto)-1])
		}
	}

	// Orquestra a higienização, tokenização e avaliação distribuída da expressão
	executarCalculo := func() {
		expressao := strings.TrimSpace(displayEntry.Text)
		if expressao == "" {
			return
		}

		exprHigienizada, err := parser.Higienizar(expressao)
		if err != nil {
			historicoLabel.SetText("Erro: expressao invalida")
			limparNoProximoInput = true
			return
		}

		tokens, err := parser.Tokenizar(exprHigienizada)
		if err != nil {
			historicoLabel.SetText(err.Error())
			limparNoProximoInput = true
			return
		}

		resultado, err := calc.AvaliarRPN(tokens, ops)
		if err != nil {
			historicoLabel.SetText(err.Error())
			limparNoProximoInput = true
			return
		}

		historicoLabel.SetText(expressao + " =")
		displayEntry.SetText(fmt.Sprintf("%g", resultado))
		limparNoProximoInput = true
	}

	// Dispara o cálculo ao pressionar Enter no teclado físico
	displayEntry.OnSubmitted = func(_ string) {
		executarCalculo()
	}

	// Grade 5x4 com os botões virtuais da calculadora
	teclado := container.NewGridWithColumns(4,
		widget.NewButton("C", limparTudo),
		widget.NewButton("(", func() { inserirTexto("(") }),
		widget.NewButton(")", func() { inserirTexto(")") }),
		widget.NewButton("+", func() { inserirTexto("+") }),

		widget.NewButton("7", func() { inserirTexto("7") }),
		widget.NewButton("8", func() { inserirTexto("8") }),
		widget.NewButton("9", func() { inserirTexto("9") }),
		widget.NewButton("-", func() { inserirTexto("-") }),

		widget.NewButton("4", func() { inserirTexto("4") }),
		widget.NewButton("5", func() { inserirTexto("5") }),
		widget.NewButton("6", func() { inserirTexto("6") }),
		widget.NewButton("*", func() { inserirTexto("*") }),

		widget.NewButton("1", func() { inserirTexto("1") }),
		widget.NewButton("2", func() { inserirTexto("2") }),
		widget.NewButton("3", func() { inserirTexto("3") }),
		widget.NewButton("/", func() { inserirTexto("/") }),

		widget.NewButton("DEL", apagarUltimo),
		widget.NewButton("0", func() { inserirTexto("0") }),
		widget.NewButton(".", func() { inserirTexto(".") }),
		widget.NewButton("=", executarCalculo),
	)

	// Painel vertical agrupando histórico e campo de entrada
	displayContainer := container.NewVBox(
		historicoLabel,
		displayEntry,
	)

	// Disposição de borda com o mostrador no topo e o teclado no centro
	conteudo := container.NewBorder(
		container.NewPadded(displayContainer),
		nil, nil, nil,
		teclado,
	)

	janela.SetContent(conteudo)
	// Ativa o foco no campo de texto
	janela.Canvas().Focus(displayEntry)
	// Exibe a janela e inicia o loop de eventos
	janela.ShowAndRun()
}
