package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func gui() {
	aplicacao := app.New()
	janela := aplicacao.NewWindow("Calculadora Distribuida")
	janela.Resize(fyne.NewSize(320, 420))

	display := widget.NewEntry()
	display.SetPlaceHolder("0")

	adicionarCaractere := func(c string) {
		display.SetText(display.Text + c)
	}

	limparDisplay := func() {
		display.SetText("")
	}

	apagarUltimo := func() {
		texto := display.Text
		if len(texto) > 0 {
			display.SetText(texto[:len(texto)-1])
		}
	}

	calcular := func() {
		if display.Text == "" {
			return
		}
		exprHigienizada, err := higienizar(display.Text)
		if err != nil {
			display.SetText("Erro: expressao invalida")
			return
		}
		tokens, err := tokenizar(exprHigienizada)
		if err != nil {
			display.SetText("Erro: tokenizacao invalida")
			return
		}
		_ = tokens
	}

	teclado := container.NewGridWithColumns(4,
		widget.NewButton("C", limparDisplay),
		widget.NewButton("(", func() { adicionarCaractere("(") }),
		widget.NewButton(")", func() { adicionarCaractere(")") }),
		widget.NewButton("/", func() { adicionarCaractere("/") }),

		widget.NewButton("7", func() { adicionarCaractere("7") }),
		widget.NewButton("8", func() { adicionarCaractere("8") }),
		widget.NewButton("9", func() { adicionarCaractere("9") }),
		widget.NewButton("*", func() { adicionarCaractere("*") }),

		widget.NewButton("4", func() { adicionarCaractere("4") }),
		widget.NewButton("5", func() { adicionarCaractere("5") }),
		widget.NewButton("6", func() { adicionarCaractere("6") }),
		widget.NewButton("-", func() { adicionarCaractere("-") }),

		widget.NewButton("1", func() { adicionarCaractere("1") }),
		widget.NewButton("2", func() { adicionarCaractere("2") }),
		widget.NewButton("3", func() { adicionarCaractere("3") }),
		widget.NewButton("+", func() { adicionarCaractere("+") }),

		widget.NewButton("0", func() { adicionarCaractere("0") }),
		widget.NewButton(".", func() { adicionarCaractere(".") }),
		widget.NewButton("DEL", apagarUltimo),
		widget.NewButton("=", calcular),
	)

	conteudo := container.NewBorder(
		container.NewVBox(display),
		nil, nil, nil,
		teclado,
	)

	janela.SetContent(conteudo)
	janela.ShowAndRun()
}
