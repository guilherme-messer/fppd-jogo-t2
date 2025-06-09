package main

import (
	"fmt"
	"sync"
	"time"
)

type Jogo struct {
	Jogadores           map[string]PosicaoJogador // a chave é o ID do jogador
	UltimoPing          map[string]time.Time
	DiamanteFoiColetado bool
	idJogadores         int
	mutex               sync.Mutex
}

type PosicaoJogador struct {
	PosX, PosY int
}

type RegistroInicial struct {
	ID      string
	Posicao PosicaoJogador
}

type Movimento struct {
	ID   string
	X, Y int
}

func jogoNovo() Jogo {
	return Jogo{
		Jogadores:           make(map[string]PosicaoJogador),
		UltimoPing:          make(map[string]time.Time),
		DiamanteFoiColetado: false,
	}
}

func (j *Jogo) SolicitarID(_ struct{}, reply *string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	j.idJogadores++
	id := fmt.Sprintf("jogador_%d", j.idJogadores)
	*reply = id
	return nil
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

func (j *Jogo) RegistrarJogador(id string, reply *PosicaoJogador) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	// Já registrado? Retorna a posição atual
	if pos, ok := j.Jogadores[id]; ok {
		*reply = pos
		return nil
	}

	// Define posição inicial (exemplo: sempre no 1,1)
	inicial := PosicaoJogador{PosX: 1, PosY: 1}
	j.Jogadores[id] = inicial
	*reply = inicial
	return nil
}

func (j *Jogo) RegistrarPosicaoInicial(req RegistroInicial, reply *bool) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if _, existe := j.Jogadores[req.ID]; existe {
		*reply = false
		return nil
	}

	j.Jogadores[req.ID] = req.Posicao
	*reply = true
	return nil
}

func (j *Jogo) GetTodosJogadores(_ struct{}, reply *map[string]PosicaoJogador) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	copia := make(map[string]PosicaoJogador)
	for id, pos := range j.Jogadores {
		copia[id] = pos
	}
	*reply = copia
	return nil
}

func (j *Jogo) AtualizarPosicao(m Movimento, reply *bool) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	// Verifica se a nova posição colide com outro jogador
	for id, pos := range j.Jogadores {
		if id != m.ID && pos.PosX == m.X && pos.PosY == m.Y {
			*reply = false // Colisão, não pode mover
			return nil
		}
	}

	// Atualiza posição no mapa
	j.Jogadores[m.ID] = PosicaoJogador{PosX: m.X, PosY: m.Y}
	*reply = true
	return nil
}

func (j *Jogo) RemoverJogador(id string, resposta *bool) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	delete(j.Jogadores, id)
	delete(j.UltimoPing, id)
	*resposta = true
	return nil
}

func (j *Jogo) Ping(id string, reply *bool) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if _, ok := j.Jogadores[id]; ok {
		j.UltimoPing[id] = time.Now()
		*reply = true
		return nil
	}
	*reply = false
	return nil
}
