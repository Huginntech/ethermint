package websocket

import (
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"

	"github.com/cosmos/ethermint/rpc/websockets"
	"github.com/cosmos/ethermint/server/config"
)

type Service struct {
	websocketServer *websockets.Server
}

// NewService creates a new gRPC server instance with a defined listener address.
func NewService(clientCtx client.Context) *Service {
	return &Service{
		websocketServer: websockets.NewServer(clientCtx),
	}
}

// Name returns the JSON-RPC service name
func (Service) Name() string {
	return "Ethereum Websocket"
}

// Start runs the websocket server
func (s Service) Start(cfg config.Config) error {
	if !cfg.EthereumWebsocket.Enable {
		return nil
	}

	s.websocketServer.Address = cfg.JSONRPC.Address

	ws := mux.NewRouter()
	ws.Handle("/", s.websocketServer)

	errCh := make(chan error)
	go func() {
		err := http.ListenAndServe(cfg.EthereumWebsocket.Address, ws)
		if err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(5 * time.Second): // assume server started successfully
		return nil
	}
}
