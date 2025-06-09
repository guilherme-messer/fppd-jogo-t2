// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool // Indica se o elemento bloqueia passagem
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa                [][]Elemento // grade 2D representando o mapa
	PosX, PosY          int          // posição atual do personagem
	UltimoVisitado      Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg           string       // mensagem para a barra de status
	Jogadores           map[string]PosicaoJogadores
	DiamanteFoiColetado bool
}

// Armazena posição dos jogadores
type PosicaoJogadores struct {
	PosX, PosY int
}

type RegistroInicial struct {
	ID      string
	Posicao PosicaoJogadores
}

// usada para encapsular o PosX, PosY do jogo, e o meuID
// o net/rpc não aceita os argumentos soltos
type Movimento struct {
	ID string
	X  int
	Y  int
}

// Elementos visuais do jogo
var (
	Personagem       = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo          = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede           = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao        = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio            = Elemento{' ', CorPadrao, CorPadrao, false}
	Diamante         = Elemento{'◆', CorAmarela, CorPadrao, true}
	DiamanteColetado = Elemento{'◇', CorAmarela, CorPadrao, true}
	OutrosJogadores  = Elemento{'☺', CorAmarela, CorPadrao, true}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{UltimoVisitado: Vazio}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Inimigo
			case Vegetacao.simbolo:
				e = Vegetacao
			case Diamante.simbolo:
				e = Diamante
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	jogo.Mapa[y][x] = jogo.UltimoVisitado   // restaura o conteúdo anterior
	jogo.UltimoVisitado = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
	jogo.Mapa[ny][nx] = elemento            // move o elemento
}

func sincronizarDiamante(jogo *Jogo, cliente *rpc.Client) {
	// Consulta o servidor para saber se o diamante já foi coletado
	var foiColetado bool
	err := cliente.Call("Jogo.GetDiamanteColetado", struct{}{}, &foiColetado)
	if err != nil {
		jogo.StatusMsg = fmt.Sprintf("Erro ao sincronizar diamante: %v", err)
		return
	}

	// Se o cliente ainda não sabia que foi coletado, atualiza o mapa
	if foiColetado && !jogo.DiamanteFoiColetado {
		jogo.DiamanteFoiColetado = true
		substituirDiamantePorColetado(jogo)
	}
}

func capturarSinalSaida(cliente *rpc.Client, meuID string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		var removido bool
		cliente.Call("Jogo.RemoverJogador", meuID, &removido)
		os.Exit(0)
	}()
}
