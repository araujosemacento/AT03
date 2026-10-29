# Calculadora Distribuída (AT03)

Calculadora com interface gráfica (Fyne) que decompõe expressões matemáticas no cliente e envia cada operação básica para ser resolvida no servidor. O projeto implementa quatro protocolos diferentes para comparar os modelos de comunicação em sistemas distribuídos.

## Protocolos Disponíveis

| Protocolo | Servidor | Cliente | Porta / Endpoint | Descrição |
| :--- | :--- | :--- | :--- | :--- |
| **Sockets TCP** | `socket-server` | `socket-client` | `:6001` | Troca de mensagens em texto puro direto pelo socket TCP |
| **RMI (Go RPC)** | `rmi-server` | `rmi-client` | `:6002` | Chamada de métodos em objeto remoto com serialização binária (`gob`) |
| **JSON-RPC** | `rpc-server` | `rpc-client` | `:6003` | Chamada remota de procedimentos com mensagens em JSON sobre TCP |
| **HTTP-REST** | `rest-server` | `rest-client` | `http://localhost:6004/calcular` | Requisições `POST` em JSON com política de Retry (3 tentativas com espera de 1s) |

## Como Funciona a Decomposição no Cliente

1. **Limpeza (`Higienizar`):** separa números, operadores e parênteses digitados pelo usuário (ex: `(10 + 15) * 4`).
2. **Ordem de Cálculo (`Tokenizar`):** usa o algoritmo Shunting-yard para converter a expressão para Notação Polonesa Reversa (`10 15 + 4 *`), respeitando parênteses, sinais negativos unários e precedência.
3. **Envio ao Servidor (`AvaliarRPN`):** o cliente lê a fila RPN e, a cada operador (`+`, `-`, `*`, `/`), faz uma requisição ao servidor pedindo o cálculo daqueles dois valores.
4. **Logs no Terminal:** o terminal do servidor mostra em tempo real cada operação individual que chega e o resultado devolvido.

---

## Forma Recomendada: Executando pelos Binários Prontos

Você não precisa instalar o Go nem compilar nada para testar o projeto. Os binários já compilados para **Linux (amd64)** e **Windows (amd64)** estão disponíveis em dois lugares:

1. **Na aba [Releases](https://github.com/araujosemacento/Calculadora-AT03/releases) do GitHub:** baixe o pacote `.tar.gz` (Linux) ou `.zip` (Windows) do protocolo desejado na última versão lançada e extraia os arquivos.
2. **Na pasta `build/` do repositório:** caso tenha clonado o projeto, os executáveis ficam organizados em `build/<protocolo>/linux/` e `build/<protocolo>/windows/`.

### Passo a passo para rodar (em qualquer protocolo)

Sempre inicie o **servidor** primeiro em um terminal (para acompanhar os logs das operações chegando) e depois abra o **cliente**.

#### No Linux

```bash
# 1. Escolha o protocolo (socket, rmi, rpc ou rest) e rode o servidor num terminal:
./socket-server

# 2. Em outro terminal (ou dando dois cliques), abra o cliente correspondente:
./socket-client
```

*(Se tiver baixado da pasta `build/`, use `./build/socket/linux/socket-server` e `./build/socket/linux/socket-client`, trocando `socket` por `rmi`, `rpc` ou `rest` conforme o teste).*

#### No Windows

```powershell
# 1. Em um terminal (PowerShell ou CMD), inicie o servidor do protocolo escolhido:
.\socket-server.exe

# 2. Dê dois cliques no cliente correspondente ou rode em outro terminal:
.\socket-client.exe
```

---

## Compilando e Rodando pelo Código-Fonte (Opcional)

Caso queira modificar o código ou compilar os executáveis na sua própria máquina, instale os pré-requisitos abaixo.

### Pré-requisitos para Compilação

- **Go (versão 1.24 ou superior)** instalado e configurado no `PATH`.
- **Compilador C (GCC) e bibliotecas gráficas** (exigidos pela biblioteca gráfica Fyne):
  - **Linux (Ubuntu / Debian):**

    ```bash
    sudo apt-get update
    sudo apt-get install -y gcc libgl1-mesa-dev xorg-dev libxkbcommon-dev libwayland-dev
    ```

  - **Linux (Fedora):**

    ```bash
    sudo dnf install -y gcc golang libXcursor-devel libXrandr-devel mesa-libGL-devel libXi-devel libXinerama-devel libXxf86vm-devel libxkbcommon-devel wayland-devel
    ```

  - **Windows:** ter um compilador GCC com suporte a 64 bits instalado, como [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) ou [MSYS2 (MinGW-w64)](https://www.msys2.org/).
- *(Opcional para compilação cruzada de clientes GUI)*: `fyne-cross` junto com **Podman** ou **Docker**.

### Rodando direto com `go run`

Abra dois terminais na raiz do projeto:

```bash
# 1. Sockets TCP (Porta :6001)
go run ./cmd/socket-server
go run ./cmd/socket-client

# 2. RMI / Go RPC (Porta :6002)
go run ./cmd/rmi-server
go run ./cmd/rmi-client

# 3. JSON-RPC (Porta :6003)
go run ./cmd/rpc-server
go run ./cmd/rpc-client

# 4. HTTP-REST com Retry (Porta :6004)
go run ./cmd/rest-server
go run ./cmd/rest-client
```

### Gerando todos os binários localmente

```bash
# No Linux
./build.sh

# No Windows
build.bat
```

Os executáveis gerados ficarão salvos dentro da pasta `build/<protocolo>/<sistema>/`.
