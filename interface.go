package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func gui() {
	aplicacao := app.New()
	janela := aplicacao.NewWindow("Calculadora")
	janela.Resize(fyne.NewSize(320, 420))

	historicoLabel := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})
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

	limparTudo := func() {
		historicoLabel.SetText("")
		displayEntry.SetText("")
		limparNoProximoInput = false
	}

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

	executarCalculo := func() {
		expressao := strings.TrimSpace(displayEntry.Text)
		if expressao == "" {
			return
		}

		exprHigienizada, err := higienizar(expressao)
		if err != nil {
			historicoLabel.SetText("Erro: expressao invalida")
			limparNoProximoInput = true
			return
		}

		tokens, err := tokenizar(exprHigienizada)
		if err != nil {
			historicoLabel.SetText(err.Error())
			limparNoProximoInput = true
			return
		}

		resultado, err := calcular(tokens)
		if err != nil {
			historicoLabel.SetText(err.Error())
			limparNoProximoInput = true
			return
		}

		historicoLabel.SetText(expressao + " =")
		displayEntry.SetText(fmt.Sprintf("%g", resultado))
		limparNoProximoInput = true
	}

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
