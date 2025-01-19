package agentprotocol_test

import (
	"fmt"
	"io"
	"testing"

	proto "go.containerssh.io/libcontainerssh/agentprotocol"

	log "go.containerssh.io/libcontainerssh/log"
)

func TestConnectionSetup(t *testing.T) {
	logger := log.NewTestLogger(t)
	fromClientReader, fromClientWriter := io.Pipe()
	toClientReader, toClientWriter := io.Pipe()

	clientCtx := proto.NewForwardCtx(toClientReader, fromClientWriter, logger)
	serverCtx := proto.NewForwardCtx(fromClientReader, toClientWriter, logger)

	closeChan := make(chan struct{})
	startedChan := make(chan struct{})

	go func() {
		connChan, err := serverCtx.StartReverseForwardClient(
			"127.0.0.1",
			8080,
			false,
		)
		if err != nil {
			panic(err)
		}

		close(startedChan)

		testConServer := <-connChan
		err = testConServer.Accept()
		if err != nil {
			logger.Error("Error accept connection", err)
		}
		buf := make([]byte, 512)
		nBytes, err := testConServer.Read(buf)
		if err != nil {
			logger.Error("Failed to read from server")
		}
		_, err = testConServer.Write(buf[:nBytes])
		if err != nil {
			logger.Error("Failed to write to server")
		}
		<-closeChan
		serverCtx.Kill()
	}()

	conType, setup, connectionChan, err := clientCtx.StartServer()
	if err != nil {
		t.Fatal("Test failed with error", err)
	}
	if conType != proto.CONNECTION_TYPE_PORT_FORWARD {
		panic(fmt.Errorf("invalid connection type %d", conType))
	}

	go func() {
		for {
			conn, ok := <-connectionChan
			if !ok {
				break
			}
			_ = conn.Reject()
			_ = conn.Close()
		}
	}()

	testConClient, err := clientCtx.NewConnectionTCP(
		setup.BindHost,
		setup.BindPort,
		"127.0.0.5",
		8081,
		func() error {
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 512)
	_, err = testConClient.Write([]byte("Message to server"))
	if err != nil {
		t.Fatal(err)
	}
	nBytes, err := testConClient.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:nBytes]) != "Message to server" {
		t.Fatalf("Expected to read 'Message to server' instead got %s", string(buf[:nBytes]))
	}
	_ = testConClient.Close()
	clientCtx.Kill()
	close(closeChan)
}
