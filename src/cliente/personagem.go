// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"fmt"
	"net/rpc"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1 // Move para cima
	case 'a':
		dx = -1 // Move para a esquerda
	case 's':
		dy = 1 // Move para baixo
	case 'd':
		dx = 1 // Move para a direita
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	// Verifica se o movimento é permitido e realiza a movimentação
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
	}
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemInteragir(jogo *Jogo, cliente *rpc.Client) {
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
				} else if sucesso {
					jogo.DiamanteFoiColetado = true
					Diamante.simbolo = '◇'
					jogo.StatusMsg = "Diamante coletado com sucesso!"
				} else {
					jogo.DiamanteFoiColetado = true
					Diamante.simbolo = '◇'
					jogo.StatusMsg = "Diamante já foi coletado."
				}
				return // Só coleta um por vez
			}
		}
	}
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo, cliente *rpc.Client) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo, cliente)
	case "mover":
		// Move o personagem com base na tecla
		personagemMover(ev.Tecla, jogo)
	}
	return true // Continua o jogo
}
