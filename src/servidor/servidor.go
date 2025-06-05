package main

type Jogo struct {
	Jogadores           map[string]PosicaoJogador // a chave é o ID do jogador
	DiamanteFoiColetado bool
}

type PosicaoJogador struct {
	PosX, PosY int
}

func jogoNovo() Jogo {
	// Inicia o diamante como não coletado
	return Jogo{DiamanteFoiColetado: false}
}

//func ColetaDiamante() {
//
//}
