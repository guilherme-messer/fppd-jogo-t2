// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"fmt"
	"net/rpc"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo, cliente *rpc.Client, meuID string) {
	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1
	case 'a':
		dx = -1
	case 's':
		dy = 1
	case 'd':
		dx = 1
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy

	// Verifica se o movimento é permitido no mapa
	if jogoPodeMoverPara(jogo, nx, ny) {
		// Verifica com o servidor se outro jogador está na posição
		mov := Movimento{ID: meuID, X: nx, Y: ny}
		var podeMover bool
		err := cliente.Call("Jogo.AtualizarPosicao", mov, &podeMover)
		if err != nil {
			jogo.StatusMsg = fmt.Sprintf("Erro RPC ao mover: %v", err)
			return
		}

		if podeMover {
			jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
			jogo.PosX, jogo.PosY = nx, ny
		} else {
			jogo.StatusMsg = "Outro jogador está nesse local."
		}
	}
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemInteragir(jogo *Jogo, cliente *rpc.Client, meuID string) {
	// Impede interação se já foi coletado
	if jogo.DiamanteFoiColetado {
		jogo.StatusMsg = "O diamante já foi coletado."
		return
	}

	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)

	// posicoes para interação com o diamante
	posicoes := [8][2]int{
		{0, -1},  // cima
		{0, 1},   // baixo
		{-1, 0},  // esquerda
		{1, 0},   // direita
		{-1, -1}, // diagonal esquerda cima
		{1, -1},  // diagonal direita cima
		{-1, 1},  // diagonal esquerda baixo
		{1, 1},   // diagonal direita baixo
	}

	for _, p := range posicoes {
		nx, ny := jogo.PosX+p[0], jogo.PosY+p[1]

		if ny >= 0 && ny < len(jogo.Mapa) && nx >= 0 && nx < len(jogo.Mapa[ny]) {
			if jogo.Mapa[ny][nx] == Diamante {
				jogo.StatusMsg = fmt.Sprintf("Baú interagido em (%d, %d)", nx, ny)

				var sucesso bool
				err := cliente.Call("Jogo.ColetarDiamante", struct{}{}, &sucesso)
				if err != nil {
					jogo.StatusMsg = fmt.Sprintf("Erro ao coletar diamante: %v", err)
				} else {
					// Diamante.simbolo = '◇'
					// ele atualiza o diamante foi coletado no for da main
					// jogo.DiamanteFoiColetado = true

					if sucesso {
						jogo.StatusMsg = "Diamante coletado com sucesso!"
					} else {
						jogo.StatusMsg = "Diamante já foi coletado."
					}

					sincronizarDiamante(jogo, cliente)

					interfaceDesenharJogo(jogo)
				}
				return
			}
		}
	}
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo, cliente *rpc.Client, meuID string) bool {
	switch ev.Tipo {
	case "sair":
		return false
	case "interagir":
		personagemInteragir(jogo, cliente, meuID)
	case "mover":
		personagemMover(ev.Tecla, jogo, cliente, meuID)
	}
	return true
}
