#!/usr/bin/env bash
set -e

# Script de compilação multiplataforma para os binários da calculadora.
# Estrutura de saída: build/<protocolo>/<sistema>/<binários>

# Lista de protocolos suportados
# Adicionar novos protocolos aqui conforme forem implementados (exemplo: "socket" "rpc" "rest" "rmi")
PROTOCOLOS=("socket")

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

echo "Iniciando compilação dos binários..."

# Assegura que o diretório bin do Go esteja presente no PATH
if command -v go &> /dev/null; then
    GOPATH_BIN="$(go env GOPATH)/bin"
    if [[ ":$PATH:" != *":$GOPATH_BIN:"* ]] && [ -d "$GOPATH_BIN" ]; then
        export PATH="$PATH:$GOPATH_BIN"
    fi
fi

# Detecta o motor de containers (prioriza Podman, com fallback para Docker)
TEM_CONTAINER=false
if command -v podman &> /dev/null; then
    export FYNE_CROSS_ENGINE="podman"
    TEM_CONTAINER=true
    echo "Motor de containers detectado: Podman"
elif command -v docker &> /dev/null && docker info &> /dev/null; then
    export FYNE_CROSS_ENGINE="docker"
    TEM_CONTAINER=true
    echo "Motor de containers detectado: Docker"
fi

for PROTO in "${PROTOCOLOS[@]}"; do
    SERVER_DIR="./cmd/${PROTO}-server"
    CLIENT_DIR="./cmd/${PROTO}-client"
    APP_ID="com.calculadora.${PROTO}"

    DIR_LINUX="build/${PROTO}/linux"
    DIR_WINDOWS="build/${PROTO}/windows"

    mkdir -p "$DIR_LINUX" "$DIR_WINDOWS"

    echo "=========================================="
    echo "Compilando protocolo: ${PROTO}"
    echo "=========================================="

    # 1. Compilação dos Servidores (Go puro: compilação cruzada direta sem Cgo)
    if [ -d "$SERVER_DIR" ]; then
        echo "[Servidor] Compilando para Linux (amd64)..."
        GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "${DIR_LINUX}/${PROTO}-server" "$SERVER_DIR"

        echo "[Servidor] Compilando para Windows (amd64)..."
        GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o "${DIR_WINDOWS}/${PROTO}-server.exe" "$SERVER_DIR"
    fi

    # 2. Compilação dos Clientes (Fyne / GUI)
    if [ -d "$CLIENT_DIR" ]; then
        # Compilação nativa para Linux
        echo "[Cliente] Compilando nativamente para Linux..."
        go build -o "${DIR_LINUX}/${PROTO}-client" "$CLIENT_DIR"

        # Compilação cruzada com fyne-cross para Windows
        if [ "$TEM_CONTAINER" = true ] && command -v fyne-cross &> /dev/null; then
            echo "[Cliente] Compilando para Windows com fyne-cross (${FYNE_CROSS_ENGINE})..."
            fyne-cross windows --engine "$FYNE_CROSS_ENGINE" -app-id "$APP_ID" -output "${PROTO}-client.exe" "$CLIENT_DIR"
            if [ -f "fyne-cross/bin/windows-amd64/${PROTO}-client.exe" ]; then
                cp "fyne-cross/bin/windows-amd64/${PROTO}-client.exe" "${DIR_WINDOWS}/"
            fi
        else
            echo "[Aviso] Podman/Docker ou fyne-cross indisponível para compilação cruzada de GUI."
            echo "[Aviso] O cliente nativo Linux e os servidores de todos os SOs foram compilados com sucesso."
        fi
    fi
done

echo "=========================================="
echo "Processo de compilação concluído."
echo "Binários disponíveis na pasta: ${ROOT_DIR}/build/"
echo "=========================================="
