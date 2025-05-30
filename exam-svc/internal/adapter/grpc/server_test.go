package service

import (
	"fmt"
	"net"
	"testing"

	"github.com/mephirious/helper-for-teachers/services/exam-svc/config"
	"github.com/mephirious/helper-for-teachers/services/exam-svc/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type mockListener struct {
	acceptChan chan net.Conn
	closeChan  chan struct{}
	addr       net.Addr
}

func newMockListener(addr string) *mockListener {
	return &mockListener{
		acceptChan: make(chan net.Conn),
		closeChan:  make(chan struct{}),
		addr:       &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8080},
	}
}

func (m *mockListener) Accept() (net.Conn, error) {
	select {
	case conn := <-m.acceptChan:
		return conn, nil
	case <-m.closeChan:
		return nil, net.ErrClosed
	}
}

func (m *mockListener) Close() error {
	close(m.closeChan)
	return nil
}

func (m *mockListener) Addr() net.Addr {
	return m.addr
}

func TestNewGRPCServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUC := mocks.NewMockTaskUseCase(ctrl)
	questionUC := mocks.NewMockQuestionUseCase(ctrl)
	examUC := mocks.NewMockExamUseCase(ctrl)

	cfg := config.Config{
		Server: config.Server{
			GRPCServer: config.GRPCServer{
				Port: 8080,
				Mode: "release",
			},
		},
	}

	t.Run("successful initialization", func(t *testing.T) {
		originalNetListen := netListen
		defer func() { netListen = originalNetListen }()
		netListen = func(network, address string) (net.Listener, error) {
			return newMockListener(address), nil
		}

		server, err := NewGRPCServer(cfg, taskUC, questionUC, examUC)
		require.NoError(t, err, "expected no error when initializing GRPCServer")
		assert.NotNil(t, server, "server should not be nil")
		assert.Equal(t, "0.0.0.0:8080", server.addr, "address should match config")
		assert.NotNil(t, server.server, "gRPC server should be initialized")
		assert.NotNil(t, server.listener, "listener should be initialized")
		assert.Equal(t, cfg.Server.GRPCServer, server.Cfg, "config should match input")
	})

	t.Run("invalid port", func(t *testing.T) {
		originalNetListen := netListen
		defer func() { netListen = originalNetListen }()
		netListen = func(network, address string) (net.Listener, error) {
			return nil, fmt.Errorf("failed to bind to port")
		}

		cfg.Server.GRPCServer.Port = -1
		server, err := NewGRPCServer(cfg, taskUC, questionUC, examUC)
		assert.Error(t, err, "expected error when binding to invalid port")
		assert.Contains(t, err.Error(), "failed to listen", "error should mention listen failure")
		assert.Nil(t, server, "server should be nil on error")
	})
}

var netListen = net.Listen
