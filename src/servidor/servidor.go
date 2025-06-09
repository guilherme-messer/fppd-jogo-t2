package main

import (
	"sync"
)

type Jogo struct {
	Jogadores           map[string]PosicaoJogador // a chave é o ID do jogador
	DiamanteFoiColetado bool
	mutex               sync.Mutex
}

type PosicaoJogador struct {
	PosX, PosY int
}

func jogoNovo() Jogo {
	// Inicia o diamante como não coletado
	return Jogo{DiamanteFoiColetado: false}
}

func (j *Jogo) ColetarDiamante(_ struct{}, reply *bool) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if j.DiamanteFoiColetado {
		*reply = false
		return nil
	}

	j.DiamanteFoiColetado = true
	*reply = true
	return nil
}

func (j *Jogo) GetDiamanteColetado(_ struct{}, reply *bool) error {
	*reply = j.DiamanteFoiColetado
	return nil
}
