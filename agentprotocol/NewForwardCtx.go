package agentprotocol

import (
	"io"

	"go.containerssh.io/libcontainerssh/log"
)

func NewForwardCtx(fromBackend io.Reader, toBackend io.Writer, logger log.Logger) *ForwardCtx {
	return &ForwardCtx{
		fromBackend: fromBackend,
		toBackend:   toBackend,
		logger:      logger,
	}
}
