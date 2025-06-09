package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func main() {
	jogo := jogoNovo()

	err := rpc.Register(&jogo)
	if err != nil {
		log.Fatalf("Erro ao registrar jogo RPC: %v", err)
	}

	// Escutar na porta 1234
	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Erro ao escutar: %v", err)
	}
	fmt.Println("Servidor ouvindo na porta 1234...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
