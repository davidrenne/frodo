// Code generated by Frodo - DO NOT EDIT.
//
//   Timestamp: Tue, 10 May 2022 16:25:33 EDT
//   Source:    games/game_service.go
//   Generator: https://github.com/davidrenne/frodo
//
package games

import (
	"context"
	"net/http"

	"github.com/davidrenne/frodo/example/multiservice/games"
	"github.com/davidrenne/frodo/rpc"
	"github.com/monadicstack/respond"
)

// NewGameServiceGateway accepts your "real" GameService instance (the thing that really does the work), and
// exposes it to other services/clients over RPC. The rpc.Gateway it returns implements http.Handler, so you
// can pass it to any standard library HTTP server of your choice.
//
//	// How to fire up your service for RPC and/or your REST API
//	service := games.GameService{ /* set up to your liking */ }
//	gateway := games.NewGameServiceGateway(service)
//	http.ListenAndServe(":8080", gateway)
//
// The default instance works well enough, but you can supply additional options such as WithMiddleware() which
// accepts any negroni-compatible middleware handlers.
func NewGameServiceGateway(service games.GameService, options ...rpc.GatewayOption) GameServiceGateway {
	gw := rpc.NewGateway(options...)
	gw.Name = "GameService"
	gw.PathPrefix = "/v2"

	gw.Register(rpc.Endpoint{
		Method:      "GET",
		Path:        "/game/:ID",
		ServiceName: "GameService",
		Name:        "GetByID",
		Handler: func(w http.ResponseWriter, req *http.Request) {
			response := respond.To(w, req)

			serviceRequest := games.GetByIDRequest{}
			if err := gw.Binder.Bind(req, &serviceRequest); err != nil {
				response.Fail(err)
				return
			}

			serviceResponse, err := service.GetByID(req.Context(), &serviceRequest)
			response.Reply(200, serviceResponse, err)
		},
	})

	gw.Register(rpc.Endpoint{
		Method:      "POST",
		Path:        "/game",
		ServiceName: "GameService",
		Name:        "Register",
		Handler: func(w http.ResponseWriter, req *http.Request) {
			response := respond.To(w, req)

			serviceRequest := games.RegisterRequest{}
			if err := gw.Binder.Bind(req, &serviceRequest); err != nil {
				response.Fail(err)
				return
			}

			serviceResponse, err := service.Register(req.Context(), &serviceRequest)
			response.Reply(201, serviceResponse, err)
		},
	})

	return GameServiceGateway{Gateway: gw, service: service}
}

type GameServiceGateway struct {
	rpc.Gateway
	service games.GameService
}

func (gw GameServiceGateway) GetByID(ctx context.Context, request *games.GetByIDRequest) (*games.GetByIDResponse, error) {
	return gw.service.GetByID(ctx, request)
}

func (gw GameServiceGateway) Register(ctx context.Context, request *games.RegisterRequest) (*games.RegisterResponse, error) {
	return gw.service.Register(ctx, request)
}

func (gw GameServiceGateway) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	gw.Gateway.ServeHTTP(w, req)
}
