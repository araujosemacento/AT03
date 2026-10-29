// Servidor HTTP-REST da calculadora.
// Fica escutando na porta :6004 e resolve as operações enviadas via requisições POST em JSON.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"at03/calculadora/pkg/calc"
)

// RequisicaoCalculo representa o JSON que chega no corpo do POST.
type RequisicaoCalculo struct {
	Operacao string  `json:"operacao"`
	A        float64 `json:"a"`
	B        float64 `json:"b"`
}

// RespostaCalculo representa o JSON devolvido para o cliente.
type RespostaCalculo struct {
	Resultado float64 `json:"resultado"`
	Erro      string  `json:"erro,omitempty"`
}

func main() {
	porta := ":6004"
	ops := calc.ServerOperations{}

	// Define a função que atende a rota /calcular
	http.HandleFunc("/calcular", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Aceita apenas requisições do tipo POST
		if r.Method != http.MethodPost {
			log.Printf("[REST] Método recusado (%s) vindo de %s", r.Method, r.RemoteAddr)
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(RespostaCalculo{Erro: "use o método POST"})
			return
		}

		// Lê o JSON enviado pelo cliente no corpo da requisição
		var req RequisicaoCalculo
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[REST] Erro ao ler JSON de %s: %v", r.RemoteAddr, err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(RespostaCalculo{Erro: "JSON inválido na requisição"})
			return
		}

		var resultado float64
		var opErr error

		log.Printf("[REST] [%s] Requisição POST recebida: %g %s %g", r.RemoteAddr, req.A, req.Operacao, req.B)

		switch req.Operacao {
		case "+":
			resultado, opErr = ops.Somar(req.A, req.B)
		case "-":
			resultado, opErr = ops.Subtrair(req.A, req.B)
		case "*":
			resultado, opErr = ops.Multiplicar(req.A, req.B)
		case "/":
			resultado, opErr = ops.Dividir(req.A, req.B)
		default:
			opErr = fmt.Errorf("operador desconhecido: %s", req.Operacao)
		}

		// Se deu erro na conta (como divisão por zero), devolve status 400
		if opErr != nil {
			log.Printf("[REST] [%s] Erro na operação: %v", r.RemoteAddr, opErr)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(RespostaCalculo{Erro: opErr.Error()})
			return
		}

		// Se deu tudo certo, devolve status 200 com o resultado em JSON
		log.Printf("[REST] [%s] Resultado enviado: %g %s %g = %g", r.RemoteAddr, req.A, req.Operacao, req.B, resultado)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RespostaCalculo{Resultado: resultado})
	})

	log.Printf("[REST] Servidor rodando e escutando em http://localhost%s/calcular", porta)
	if err := http.ListenAndServe(porta, nil); err != nil {
		log.Fatalf("[REST] Erro ao iniciar o servidor HTTP: %v", err)
	}
}
