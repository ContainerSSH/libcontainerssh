package agentprotocol

import (
	"io"

	log "go.containerssh.io/libcontainerssh/log"
)

func NewForwardCtx(fromBackend io.Reader, toBackend io.Writer, logger log.Logger) ForwardCtx {
	return &forwardCtx{
		fromBackend: fromBackend,
		toBackend:   toBackend,
		logger:      logger,
	}
}
