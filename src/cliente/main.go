package main

import (
	"log"
	"net/rpc"
	"os"
)

func main() {
	// Conectar ao servidor
	cliente, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("Erro ao conectar ao servidor RPC: %v", err)
	}

	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo local
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo, cliente); !continuar {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}
