package proxyproto_test

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/containerssh/libcontainerssh/internal/proxyproto"
	goproxyproto "github.com/pires/go-proxyproto"
)

type fakeConn struct {
	remoteAddr string
	localAddr  string
	pipeReader io.ReadCloser
	pipeWriter io.WriteCloser
}

func NewFakeConn(clientAddr string, serverAddr string) (fakeConn, fakeConn) {
	clientPipeReader, clientPipeWriter := io.Pipe()
	serverPipeReader, serverPipeWriter := io.Pipe()
	return fakeConn{
			remoteAddr: clientAddr,
			localAddr:  serverAddr,
			pipeReader: serverPipeReader,
			pipeWriter: clientPipeWriter,
		}, fakeConn{
			remoteAddr: serverAddr,
			localAddr:  clientAddr,
			pipeReader: clientPipeReader,
			pipeWriter: serverPipeWriter,
		}
}

func (f fakeConn) Read(b []byte) (n int, err error) {
	return f.pipeReader.Read(b)
}

func (f fakeConn) Write(b []byte) (n int, err error) {
	return f.pipeWriter.Write(b)
}

func (f fakeConn) Close() error {
	f.pipeWriter.Close()
	f.pipeReader.Close()
	return nil
}

func (f fakeConn) LocalAddr() net.Addr {
	return &net.TCPAddr{
		IP: net.ParseIP(f.localAddr),
	}
}

func (f fakeConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP: net.ParseIP(f.remoteAddr),
	}
}
func (f fakeConn) SetDeadline(t time.Time) error {
	return fmt.Errorf("Unimplemented")
}
func (f fakeConn) SetReadDeadline(t time.Time) error {
	return fmt.Errorf("Unimplemented")
}
func (f fakeConn) SetWriteDeadline(t time.Time) error {
	return fmt.Errorf("Unimplemented")
}

func TestProxyWithHeader(t *testing.T) {
	clientIP := "127.0.0.1"
	proxyIP := "127.0.0.2"
	serverIP := "127.0.0.3"

	server, proxy := NewFakeConn(proxyIP, serverIP)
	wrappedConn, proxyAddr, err := proxyproto.WrapProxy(server, []string{proxyIP})
	if err != nil {
		t.Fatal(err)
	}

	header := &goproxyproto.Header{
		Version:           1,
		Command:           goproxyproto.PROXY,
		TransportProtocol: goproxyproto.TCPv4,
		SourceAddr: &net.TCPAddr{
			IP:   net.ParseIP(clientIP),
			Port: 1000,
		},
		DestinationAddr: &net.TCPAddr{
			IP:   net.ParseIP(proxyIP),
			Port: 2000,
		},
	}
	go func() {
		_, err := header.WriteTo(proxy)
		if err != nil {
			return
		}
	}()

	if proxyAddr == nil {
		t.Fatalf("Proxy info was rejected")
	}
	if proxyAddr.String() != proxyIP+":0" {
		t.Fatalf("Unexpected proxy address %s, expected %s", proxyAddr, proxyIP)
	}
	if wrappedConn.RemoteAddr().String() != clientIP+":1000" {
		t.Fatalf("Header not accepted when it should be %s != %s", wrappedConn.RemoteAddr().String(), clientIP+":1000")
	}
}

func TestProxyUnauthorizedHeader(t *testing.T) {
	clientIP := "127.0.0.1"
	proxyIP := "127.0.0.2"
	serverIP := "127.0.0.3"

	server, proxy := NewFakeConn(proxyIP, serverIP)
	_, proxyAddr, err := proxyproto.WrapProxy(server, []string{"128.0.0.2"})
	if err != nil {
		t.Fatal(err)
	}

	header := &goproxyproto.Header{
		Version:           1,
		Command:           goproxyproto.PROXY,
		TransportProtocol: goproxyproto.TCPv4,
		SourceAddr: &net.TCPAddr{
			IP:   net.ParseIP(clientIP),
			Port: 1000,
		},
		DestinationAddr: &net.TCPAddr{
			IP:   net.ParseIP(proxyIP),
			Port: 2000,
		},
	}
	go func() {
		_, err := header.WriteTo(proxy)
		if err != nil {
			return
		}
	}()

	if proxyAddr != nil {
		t.Fatalf("Proxy info was accepted when unauthorized")
	}
}
