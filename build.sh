#!/usr/bin/env bash
set -e

# Script de compilação multiplataforma para os binários da calculadora.
# Estrutura de saída: build/<protocolo>/<sistema>/<binários>

# Lista de protocolos suportados
PROTOCOLOS=("socket" "rmi" "rpc" "rest")

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

echo "Iniciando compilação dos binários..."

# Garante que a pasta bin do Go esteja no PATH
if command -v go &> /dev/null; then
    GOPATH_BIN="$(go env GOPATH)/bin"
    if [[ ":$PATH:" != *":$GOPATH_BIN:"* ]] && [ -d "$GOPATH_BIN" ]; then
        export PATH="$PATH:$GOPATH_BIN"
    fi
fi

# Verifica se há Podman ou Docker instalado para compilar a interface de Windows a partir do Linux
TEM_CONTAINER=false
if command -v podman &> /dev/null; then
    export FYNE_CROSS_ENGINE="podman"
    TEM_CONTAINER=true
    echo "Motor de containers encontrado: Podman"
elif command -v docker &> /dev/null && docker info &> /dev/null; then
    export FYNE_CROSS_ENGINE="docker"
    TEM_CONTAINER=true
    echo "Motor de containers encontrado: Docker"
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

    # 1. Compilação dos Servidores (Go puro, sem precisar de Cgo)
    if [ -d "$SERVER_DIR" ]; then
        echo "[Servidor] Compilando para Linux (amd64)..."
        GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o "${DIR_LINUX}/${PROTO}-server" "$SERVER_DIR"

        echo "[Servidor] Compilando para Windows (amd64)..."
        GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o "${DIR_WINDOWS}/${PROTO}-server.exe" "$SERVER_DIR"
    fi

    # 2. Compilação dos Clientes (interface gráfica Fyne)
    if [ -d "$CLIENT_DIR" ]; then
        echo "[Cliente] Compilando nativamente para Linux..."
        go build -ldflags="-s -w" -o "${DIR_LINUX}/${PROTO}-client" "$CLIENT_DIR"

        if [ "$TEM_CONTAINER" = true ] && command -v fyne-cross &> /dev/null; then
            echo "[Cliente] Compilando para Windows com fyne-cross (${FYNE_CROSS_ENGINE})..."
            fyne-cross windows --engine "$FYNE_CROSS_ENGINE" -app-id "$APP_ID" -output "${PROTO}-client.exe" "$CLIENT_DIR"
            if [ -f "fyne-cross/bin/windows-amd64/${PROTO}-client.exe" ]; then
                cp "fyne-cross/bin/windows-amd64/${PROTO}-client.exe" "${DIR_WINDOWS}/"
            fi
        else
            echo "[Aviso] Podman/Docker ou fyne-cross não encontrado para gerar o cliente Windows."
            echo "[Aviso] O cliente Linux e os servidores de ambos os sistemas foram compilados normalmente."
        fi
    fi
done

echo "=========================================="
echo "Compilação concluída."
echo "Binários salvos em: ${ROOT_DIR}/build/"
echo "=========================================="
