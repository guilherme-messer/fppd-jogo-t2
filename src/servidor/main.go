package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

func main() {
	jogo := jogoNovo()

	err := rpc.Register(&jogo)
	if err != nil {
		log.Fatalf("Erro ao registrar jogo RPC: %v", err)
	}

	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Erro ao escutar: %v", err)
	}
	fmt.Println("Servidor ouvindo na porta 1234...")

	// Go routine para remover clientes inativos, que não fizeram RPC
	// na func Ping há mais de 10 segundos
	go func() {
		for {
			jogo.mutex.Lock()
			for id, ultimo := range jogo.UltimoPing {
				if time.Since(ultimo) > 10*time.Second {
					fmt.Printf("Jogador %s removido por inatividade\n", id)
					delete(jogo.Jogadores, id)
					delete(jogo.UltimoPing, id)
				}
			}
			jogo.mutex.Unlock()
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
