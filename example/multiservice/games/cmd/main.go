package main

import (
	"net/http"

	"github.com/davidrenne/frodo/example/multiservice/games"
	gamesrpc "github.com/davidrenne/frodo/example/multiservice/games/gen"
)

func main() {
	serviceHandler := games.GameServiceHandler{
		Repo: games.NewRepo(),
	}
	gateway := gamesrpc.NewGameServiceGateway(&serviceHandler)
	http.ListenAndServe(":9001", gateway)
}
