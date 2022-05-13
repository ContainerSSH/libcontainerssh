package proxyproto

import (
	"net"

	"github.com/pires/go-proxyproto"
)

// WrapProxy is a function that wraps a net.Conn around the PROXY tcp protocol. It is used for correctly reporting the originator IP address when a service is running behind a load balancer
// In case proxy use is allowed the wrapped network connection is returned along with the IP address of the proxy that it is used. The wrapped network connection will return the IP address
// of the client when RemoteAddr() is called
//
// conn is the network connection to wrap
// proxyList is a list of addresses that are allowed to send proxy information
//
func WrapProxy(conn net.Conn, proxyList []string) (net.Conn, *net.TCPAddr, error) {
	if len(proxyList) == 0 {
		return conn, nil, nil
	}
	policyFunc := proxyproto.MustStrictWhiteListPolicy(proxyList)
	policy, err := policyFunc(conn.RemoteAddr())
	if err != nil {
		return nil, nil, err
	}
	if policy == proxyproto.REJECT || policy == proxyproto.IGNORE {
		// If it's not an approved proxy we should fail loudly, not silently
		return conn, nil, nil
	}
	tcpAddr := conn.RemoteAddr().(*net.TCPAddr)
	return proxyproto.NewConn(
		conn,
		proxyproto.WithPolicy(policy),
	), tcpAddr, nil
}
