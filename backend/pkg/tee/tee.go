package tee

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"sync"
)

type TeeHandleFunc func(ctx context.Context, data []byte) error

var bufferPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 4096)
		return &buf
	},
}

type Tee struct {
	ctx    context.Context
	logger *slog.Logger
	Reader io.Reader
	Writer io.Writer
	ch     chan []byte
	handle TeeHandleFunc
}

func NewTee(
	ctx context.Context,
	logger *slog.Logger,
	reader io.Reader,
	writer io.Writer,
	handle TeeHandleFunc,
) *Tee {
	t := &Tee{
		ctx:    ctx,
		logger: logger,
		Reader: reader,
		Writer: writer,
		handle: handle,
		ch:     make(chan []byte, 32*1024),
	}
	go t.Handle()
	return t
}

func (t *Tee) Close() {
	select {
	case <-t.ch:
		// channel 已经关闭
	default:
		close(t.ch)
	}
}

func (t *Tee) Handle() {
	for {
		select {
		case data, ok := <-t.ch:
			if !ok {
				t.logger.DebugContext(t.ctx, "Tee Handle closed")
				return
			}
			err := t.handle(t.ctx, data)
			if err != nil {
				t.logger.With("data", string(data)).With("error", err).ErrorContext(t.ctx, "Tee Handle error")
				return
			}
		case <-t.ctx.Done():
			t.logger.DebugContext(t.ctx, "Tee Handle ctx done")
			return
		}
	}
}

func (t *Tee) Stream() error {
	bufPtr := bufferPool.Get().(*[]byte)
	buf := *bufPtr
	defer bufferPool.Put(bufPtr)

	for {
		n, err := t.Reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if n > 0 {
			_, err = t.Writer.Write(buf[:n])
			if err != nil {
				return err
			}
			if flusher, ok := t.Writer.(http.Flusher); ok {
				flusher.Flush()
			}
			data := make([]byte, n)
			copy(data, buf[:n])
			t.ch <- data
		}
	}
}
