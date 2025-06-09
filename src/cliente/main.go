package main

import (
	"log"
	"net/rpc"
	"os"
	"sync"
	"time"
)

func main() {
	// Conectar ao servidor
	cliente, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("Erro ao conectar ao servidor RPC: %v", err)
	}

	// Solicitar ID ao servidor
	var meuID string
	err = cliente.Call("Jogo.SolicitarID", struct{}{}, &meuID)
	if err != nil {
		log.Fatalf("Erro ao solicitar ID do jogador: %v", err)
	}
	log.Printf("ID atribuído: %s\n", meuID)

	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "./src/cliente/mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo local
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Registrar a posição inicial do jogador no servidor
	var registrado bool
	err = cliente.Call("Jogo.RegistrarPosicaoInicial", RegistroInicial{
		ID: meuID,
		Posicao: PosicaoJogadores{
			PosX: jogo.PosX,
			PosY: jogo.PosY,
		},
	}, &registrado)
	if err != nil || !registrado {
		log.Fatalf("Erro ao registrar posição inicial do jogador: %v", err)
	}

	// Canal para encerrar goroutines ao sair
	done := make(chan struct{})
	var fecharOnce sync.Once

	// Goroutine para sincronizar jogadores e diamante periodicamente
	go func() {
		ultimoPing := time.Now()

		for {
			select {
			case <-done:
				return
			default:
				// Atualiza jogadores
				var jogadores map[string]PosicaoJogadores
				err := cliente.Call("Jogo.GetTodosJogadores", struct{}{}, &jogadores)
				if err == nil {
					delete(jogadores, meuID)
					jogo.Jogadores = jogadores
				}

				// Atualiza diamante
				sincronizarDiamante(&jogo, cliente)

				// Atualiza interface
				interfaceDesenharJogo(&jogo)

				// Faz o ping para o servidor
				pingServidor(cliente, meuID, &ultimoPing, done, &fecharOnce)

				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo, cliente, meuID); !continuar {
			var removido bool
			_ = cliente.Call("Jogo.RemoverJogador", meuID, &removido)
			fecharOnce.Do(func() {
				close(done)
			})
			break
		}
	}
}
