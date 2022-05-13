package binary

import (
	"compress/gzip"
	"fmt"
	"net"

	"github.com/containerssh/libcontainerssh/auditlog/message"
	"github.com/containerssh/libcontainerssh/internal/auditlog/codec"
	"github.com/containerssh/libcontainerssh/internal/auditlog/storage"
	"github.com/containerssh/libcontainerssh/internal/geoip/geoipprovider"
	"github.com/fxamacker/cbor"
)

// NewEncoder creates an encoder that encodes messages in CBOR+GZIP format as documented
//            on https://containerssh.github.io/advanced/audit/format/
func NewEncoder(geoIPProvider geoipprovider.LookupProvider) codec.Encoder {
	return &encoder{
		geoIPProvider: geoIPProvider,
	}
}

type encoder struct {
	geoIPProvider geoipprovider.LookupProvider
}

func (e *encoder) GetMimeType() string {
	return "application/octet-stream"
}

func (e *encoder) GetFileExtension() string {
	return ""
}

func (e *encoder) Encode(messages <-chan message.Message, storage storage.Writer) error {
	header := newHeader(CurrentVersion).getBytes()
	if _, err := storage.Write(header); err != nil {
		return err
	}

	var gzipHandle *gzip.Writer
	var encoder *cbor.Encoder
	gzipHandle = gzip.NewWriter(storage)
	encoder = cbor.NewEncoder(gzipHandle, cbor.EncOptions{})
	if err := encoder.StartIndefiniteArray(); err != nil {
		return fmt.Errorf("failed to start infinite array (%w)", err)
	}

	startTime := int64(0)
	var ip = ""
	var proxy *string
	var country = "XX"
	var username *string
	for {
		msg, ok := <-messages
		if !ok {
			break
		}
		if startTime == 0 {
			startTime = msg.Timestamp
		}
		ip, proxy, country, username = e.storeMetadata(msg, storage, startTime, ip, proxy, country, username)
		if err := encoder.Encode(&msg); err != nil {
			return fmt.Errorf("failed to encode audit log message (%w)", err)
		}
		if msg.MessageType == message.TypeDisconnect {
			break
		}
	}
	if err := encoder.EndIndefinite(); err != nil {
		return fmt.Errorf("failed to end audit log infinite array (%w)", err)
	}
	if err := gzipHandle.Flush(); err != nil {
		return fmt.Errorf("failed to flush audit log gzip stream (%w)", err)
	}
	if err := storage.Close(); err != nil {
		return fmt.Errorf("failed to close audit log gzip stream (%w)", err)
	}
	return nil
}

func (e *encoder) storeMetadata(
	msg message.Message,
	storage storage.Writer,
	startTime int64,
	ip string,
	proxy *string,
	country string,
	username *string,
) (string, *string, string, *string) {
	switch msg.MessageType {
	case message.TypeConnect:
		payload := msg.Payload.(message.PayloadConnect)
		ip = payload.RemoteAddr
		if payload.ProxyAddr != "" {
			proxy = &payload.ProxyAddr
		}
		country := e.geoIPProvider.Lookup(net.ParseIP(ip))
		storage.SetMetadata(startTime/1000000000, ip, proxy, country, username)
	case message.TypeAuthPasswordSuccessful:
		u := msg.Payload.(message.PayloadAuthPassword).Username
		username = &u
		storage.SetMetadata(startTime/1000000000, ip, proxy, country, username)
	case message.TypeAuthPubKeySuccessful:
		payload := msg.Payload.(message.PayloadAuthPubKey)
		username = &payload.Username
		storage.SetMetadata(startTime/1000000000, ip, proxy, country, username)
	case message.TypeHandshakeSuccessful:
		payload := msg.Payload.(message.PayloadHandshakeSuccessful)
		username = &payload.Username
		storage.SetMetadata(startTime/1000000000, ip, proxy, country, username)
	}

	return ip, proxy, country, username
}
