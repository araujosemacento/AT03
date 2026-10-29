// Cliente HTTP-REST da calculadora.
// Envia requisições POST em JSON com política de Retry (espera e tenta de novo) e exibe a interface gráfica.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"at03/calculadora/pkg/ui"
)

// RequisicaoCalculo representa o JSON enviado para o servidor REST.
type RequisicaoCalculo struct {
	Operacao string  `json:"operacao"`
	A        float64 `json:"a"`
	B        float64 `json:"b"`
}

// RespostaCalculo representa o JSON que volta do servidor.
type RespostaCalculo struct {
	Resultado float64 `json:"resultado"`
	Erro      string  `json:"erro,omitempty"`
}

// RESTClient implementa [calc.Operacoes] fazendo chamadas HTTP com política de Retry.
type RESTClient struct {
	url           string
	httpClient    *http.Client
	maxTentativas int
	tempoEspera   time.Duration
}

// NovoRESTClient configura o cliente HTTP com timeout, limite de tentativas e tempo de espera.
func NovoRESTClient(url string) *RESTClient {
	return &RESTClient{
		url: url,
		httpClient: &http.Client{
			Timeout: 4 * time.Second,
		},
		maxTentativas: 3,
		tempoEspera:   1 * time.Second,
	}
}

// enviarComRetry manda o POST para o servidor e tenta de novo se houver falha de rede ou erro 5xx.
func (c *RESTClient) enviarComRetry(operacao string, a, b float64) (float64, error) {
	dados := RequisicaoCalculo{
		Operacao: operacao,
		A:        a,
		B:        b,
	}

	// Transforma a estrutura da operação em bytes JSON
	payload, err := json.Marshal(dados)
	if err != nil {
		return 0, fmt.Errorf("erro ao montar JSON: %w", err)
	}

	// Tenta enviar a requisição até o limite configurado em maxTentativas
	for tentativa := 1; tentativa <= c.maxTentativas; tentativa++ {
		req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(payload))
		if err != nil {
			return 0, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Envia a requisição HTTP para o servidor
		resp, err := c.httpClient.Do(req)

		// Se o servidor estiver fora do ar ou der erro de rede, espera e tenta de novo
		if err != nil {
			log.Printf("[REST] [Tentativa %d/%d] Falha ao falar com o servidor: %v", tentativa, c.maxTentativas, err)
			if tentativa < c.maxTentativas {
				log.Printf("[REST] Esperando %v antes de tentar de novo...", c.tempoEspera)
				time.Sleep(c.tempoEspera)
				continue
			}
			return 0, fmt.Errorf("servidor sem resposta após %d tentativas: %w", c.maxTentativas, err)
		}

		// Lê o JSON da resposta e fecha o corpo da requisição
		var resposta RespostaCalculo
		decErr := json.NewDecoder(resp.Body).Decode(&resposta)
		resp.Body.Close()

		if decErr != nil {
			log.Printf("[REST] [Tentativa %d/%d] Resposta JSON inválida: %v", tentativa, c.maxTentativas, decErr)
			if tentativa < c.maxTentativas {
				time.Sleep(c.tempoEspera)
				continue
			}
			return 0, fmt.Errorf("resposta inválida do servidor: %w", decErr)
		}

		// Erro 400 (como divisão por zero) é erro na conta em si, então não adianta tentar de novo
		if resp.StatusCode == http.StatusBadRequest {
			return 0, fmt.Errorf("%s", resposta.Erro)
		}

		// Se o servidor retornar erro interno (500+), espera e tenta de novo
		if resp.StatusCode >= 500 {
			log.Printf("[REST] [Tentativa %d/%d] Servidor respondeu com erro HTTP %d: %s", tentativa, c.maxTentativas, resp.StatusCode, resposta.Erro)
			if tentativa < c.maxTentativas {
				time.Sleep(c.tempoEspera)
				continue
			}
			return 0, fmt.Errorf("erro no servidor (HTTP %d): %s", resp.StatusCode, resposta.Erro)
		}

		// Qualquer código inesperado fora de 200 OK
		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("status HTTP inesperado (%d): %s", resp.StatusCode, resposta.Erro)
		}

		// Deu tudo certo (HTTP 200)
		return resposta.Resultado, nil
	}

	return 0, fmt.Errorf("não foi possível concluir a operação após %d tentativas", c.maxTentativas)
}

func (c *RESTClient) Somar(a, b float64) (float64, error) {
	return c.enviarComRetry("+", a, b)
}

func (c *RESTClient) Subtrair(a, b float64) (float64, error) {
	return c.enviarComRetry("-", a, b)
}

func (c *RESTClient) Multiplicar(a, b float64) (float64, error) {
	return c.enviarComRetry("*", a, b)
}

func (c *RESTClient) Dividir(a, b float64) (float64, error) {
	return c.enviarComRetry("/", a, b)
}

func main() {
	urlServidor := "http://localhost:6004/calcular"
	log.Printf("[REST] Cliente configurado para %s (retry: 3 tentativas com espera de 1s)", urlServidor)

	cliente := NovoRESTClient(urlServidor)

	log.Println("[REST] Abrindo a interface...")
	ui.IniciarGUI("Calculadora: HTTP-REST", cliente)
}
