// Pacote ui monta a janela e os botões da interface gráfica da calculadora.
package ui

import (
	"strconv"
	"strings"

	"at03/calculadora/pkg/calc"
	"at03/calculadora/pkg/parser"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// IniciarGUI cria a janela no Fyne, conecta os botões e o teclado físico,
// decompõe a expressão digitada e envia cada conta para o cliente de rede em [calc.Operacoes].
func IniciarGUI(titulo string, ops calc.Operacoes) {
	aplicacao := app.New()
	janela := aplicacao.NewWindow(titulo)
	janela.Resize(fyne.NewSize(320, 420))

	// Texto superior que mostra a expressão anterior ou avisos de erro
	historicoLabel := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})

	// Campo onde o usuário digita ou vê o resultado
	displayEntry := widget.NewEntry()
	displayEntry.SetPlaceHolder("0")

	inserirTexto := func(s string) {
		displayEntry.SetText(displayEntry.Text + s)
	}

	limparTudo := func() {
		historicoLabel.SetText("")
		displayEntry.SetText("")
	}

	apagarUltimo := func() {
		texto := displayEntry.Text
		if len(texto) > 0 {
			displayEntry.SetText(texto[:len(texto)-1])
		}
	}

	// Faz a limpeza, converte para RPN e resolve operação por operação no servidor
	executarCalculo := func() {
		expressao := strings.TrimSpace(displayEntry.Text)
		if expressao == "" {
			return
		}

		exprHigienizada, err := parser.Higienizar(expressao)
		if err != nil {
			historicoLabel.SetText("Erro: expressão inválida")
			return
		}

		tokens, err := parser.Tokenizar(exprHigienizada)
		if err != nil {
			historicoLabel.SetText(err.Error())
			return
		}

		resultado, err := calc.AvaliarRPN(tokens, ops)
		if err != nil {
			historicoLabel.SetText(err.Error())
			return
		}

		historicoLabel.SetText(expressao + " =")
		displayEntry.SetText(strconv.FormatFloat(resultado, 'f', -1, 64))
	}

	// Calcula ao apertar Enter no teclado
	displayEntry.OnSubmitted = func(_ string) {
		executarCalculo()
	}

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

	displayContainer := container.NewVBox(
		historicoLabel,
		displayEntry,
	)

	conteudo := container.NewBorder(
		container.NewPadded(displayContainer),
		nil, nil, nil,
		teclado,
	)

	janela.SetContent(conteudo)
	janela.Canvas().Focus(displayEntry)
	janela.ShowAndRun()
}
