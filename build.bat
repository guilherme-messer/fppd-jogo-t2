@echo off
REM build.bat - Compila cliente.exe e servidor.exe no Windows

REM Inicializa go.mod se necessário
IF NOT EXIST go.mod (
    echo Inicializando módulo Go...
    go mod init fppd-jogo-t2
    go get -u github.com/nsf/termbox-go
)

REM Compila cliente
echo Compilando cliente...
go build -o cliente.exe .\src\cliente

REM Compila servidor
echo Compilando servidor...
go build -o servidor.exe .\src\servidor

echo Build finalizado com sucesso.
exit /b
