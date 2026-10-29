@echo off
setlocal enabledelayedexpansion

rem Script de compilacao multiplataforma para os binarios da calculadora (Windows CMD).
rem Estrutura de saida: build\<protocolo>\<sistema>\<binarios>

rem Lista de protocolos suportados
set PROTOCOLOS=socket rmi rpc rest

set ROOT_DIR=%~dp0
cd /d "%ROOT_DIR%"

echo Iniciando compilacao dos binarios...

for %%P in (%PROTOCOLOS%) do (
    set PROTO=%%P
    set SERVER_DIR=.\cmd\%%P-server
    set CLIENT_DIR=.\cmd\%%P-client
    set APP_ID=com.calculadora.%%P

    set DIR_LINUX=build\%%P\linux
    set DIR_WINDOWS=build\%%P\windows

    if not exist "!DIR_LINUX!" mkdir "!DIR_LINUX!"
    if not exist "!DIR_WINDOWS!" mkdir "!DIR_WINDOWS!"

    echo ==========================================
    echo Compilando protocolo: %%P
    echo ==========================================

    rem 1. Compilacao dos Servidores (Go puro: sem Cgo)
    if exist "!SERVER_DIR!" (
        echo [Servidor] Compilando para Windows (amd64)...
        set GOOS=windows
        set GOARCH=amd64
        set CGO_ENABLED=0
        go build -ldflags="-s -w" -o "!DIR_WINDOWS!\%%P-server.exe" "!SERVER_DIR!"

        echo [Servidor] Compilando para Linux (amd64)...
        set GOOS=linux
        set GOARCH=amd64
        set CGO_ENABLED=0
        go build -ldflags="-s -w" -o "!DIR_LINUX!\%%P-server" "!SERVER_DIR!"
    )

    rem 2. Compilacao dos Clientes (Fyne / GUI)
    if exist "!CLIENT_DIR!" (
        echo [Cliente] Compilando nativamente para Windows...
        set GOOS=windows
        set GOARCH=amd64
        go build -ldflags="-s -w -H=windowsgui" -o "!DIR_WINDOWS!\%%P-client.exe" "!CLIENT_DIR!"

        rem Verificacao do motor de containers e fyne-cross para gerar Linux
        where fyne-cross >nul 2>nul
        if !errorlevel! equ 0 (
            set TEM_CONTAINER=0
            where podman >nul 2>nul
            if !errorlevel! equ 0 (
                set FYNE_CROSS_ENGINE=podman
                set TEM_CONTAINER=1
                echo Motor de containers detectado: Podman
            ) else (
                where docker >nul 2>nul
                if !errorlevel! equ 0 (
                    docker info >nul 2>nul
                    if !errorlevel! equ 0 (
                        set FYNE_CROSS_ENGINE=docker
                        set TEM_CONTAINER=1
                        echo Motor de containers detectado: Docker
                    )
                )
            )

            if !TEM_CONTAINER! equ 1 (
                echo [Cliente] Compilando para Linux com fyne-cross (!FYNE_CROSS_ENGINE!)...
                fyne-cross linux --engine !FYNE_CROSS_ENGINE! -app-id "!APP_ID!" -output "%%P-client" "!CLIENT_DIR!"
                if exist "fyne-cross\bin\linux-amd64\%%P-client" (
                    copy /y "fyne-cross\bin\linux-amd64\%%P-client" "!DIR_LINUX!\" >nul
                )
            ) else (
                echo [Aviso] Podman ou Docker nao detectado. Cliente Linux via fyne-cross ignorado.
            )
        ) else (
            echo [Aviso] fyne-cross nao encontrado no PATH do Windows.
        )
    )
)

echo ==========================================
echo Processo de compilacao concluido.
echo Binarios disponiveis na pasta: %ROOT_DIR%build\
echo ==========================================
