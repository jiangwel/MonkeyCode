package tee

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockWriter 模拟 Writer 接口
type mockWriter struct {
	buf     *bytes.Buffer
	delay   time.Duration
	errorOn int // 在第几次写入时返回错误
	count   int
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		buf: &bytes.Buffer{},
	}
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	m.count++
	if m.errorOn > 0 && m.count >= m.errorOn {
		return 0, errors.New("mock write error")
	}
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	return m.buf.Write(p)
}

func (m *mockWriter) String() string {
	return m.buf.String()
}

// mockReader 模拟 Reader 接口
type mockReader struct {
	data    []byte
	pos     int
	chunk   int // 每次读取的字节数
	errorOn int // 在第几次读取时返回错误
	count   int
}

func newMockReader(data string, chunk int) *mockReader {
	return &mockReader{
		data:  []byte(data),
		chunk: chunk,
	}
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	m.count++
	if m.errorOn > 0 && m.count >= m.errorOn {
		return 0, errors.New("mock read error")
	}

	if m.pos >= len(m.data) {
		return 0, io.EOF
	}

	readSize := m.chunk
	if readSize <= 0 || readSize > len(p) {
		readSize = len(p)
	}

	remaining := len(m.data) - m.pos
	if readSize > remaining {
		readSize = remaining
	}

	copy(p, m.data[m.pos:m.pos+readSize])
	m.pos += readSize
	return readSize, nil
}

// TestTeeBasicFunctionality 测试基本功能
func TestTeeBasicFunctionality(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	testData := "Hello, World! This is a test message."
	reader := newMockReader(testData, 10) // 每次读取10字节
	writer := newMockWriter()

	var handledData [][]byte
	var mu sync.Mutex

	handle := func(ctx context.Context, data []byte) error {
		mu.Lock()
		defer mu.Unlock()
		// 复制数据，因为原始数据可能被重用
		dataCopy := make([]byte, len(data))
		copy(dataCopy, data)
		handledData = append(handledData, dataCopy)
		return nil
	}

	tee := NewTee(ctx, logger, reader, writer, handle)
	defer tee.Close()

	err := tee.Stream()
	if err != nil {
		t.Fatalf("Stream() failed: %v", err)
	}

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 验证写入的数据
	if writer.String() != testData {
		t.Errorf("Expected writer data %q, got %q", testData, writer.String())
	}

	// 验证处理的数据
	mu.Lock()
	var totalHandled []byte
	for _, chunk := range handledData {
		totalHandled = append(totalHandled, chunk...)
	}
	mu.Unlock()

	if string(totalHandled) != testData {
		t.Errorf("Expected handled data %q, got %q", testData, string(totalHandled))
	}
}

// TestTeeWithErrors 测试错误处理
func TestTeeWithErrors(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("ReaderError", func(t *testing.T) {
		reader := newMockReader("test data", 5)
		reader.errorOn = 2 // 第二次读取时出错
		writer := newMockWriter()

		handle := func(ctx context.Context, data []byte) error {
			return nil
		}

		tee := NewTee(ctx, logger, reader, writer, handle)
		defer tee.Close()

		err := tee.Stream()
		if err == nil {
			t.Error("Expected error from reader, got nil")
		}
	})

	t.Run("WriterError", func(t *testing.T) {
		reader := newMockReader("test data", 5)
		writer := newMockWriter()
		writer.errorOn = 2 // 第二次写入时出错

		handle := func(ctx context.Context, data []byte) error {
			return nil
		}

		tee := NewTee(ctx, logger, reader, writer, handle)
		defer tee.Close()

		err := tee.Stream()
		if err == nil {
			t.Error("Expected error from writer, got nil")
		}
	})

	t.Run("HandleError", func(t *testing.T) {
		reader := newMockReader("test data", 5)
		writer := newMockWriter()

		handle := func(ctx context.Context, data []byte) error {
			return errors.New("handle error")
		}

		tee := NewTee(ctx, logger, reader, writer, handle)
		defer tee.Close()

		// 启动 Stream 在单独的 goroutine 中
		go func() {
			tee.Stream()
		}()

		// 等待一段时间让处理器有机会处理数据并出错
		time.Sleep(200 * time.Millisecond)
	})
}

// TestTeeContextCancellation 测试上下文取消
func TestTeeContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	// 创建一个会持续产生数据的 reader
	reader := strings.NewReader(strings.Repeat("test data ", 1000))
	writer := newMockWriter()

	handle := func(ctx context.Context, data []byte) error {
		return nil
	}

	tee := NewTee(ctx, logger, reader, writer, handle)
	defer tee.Close()

	// 在单独的 goroutine 中启动 Stream
	done := make(chan error, 1)
	go func() {
		done <- tee.Stream()
	}()

	// 等待一段时间后取消上下文
	time.Sleep(50 * time.Millisecond)
	cancel()

	// 等待 Stream 完成
	select {
	case err := <-done:
		if err != nil && err != io.EOF {
			t.Logf("Stream completed with error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Stream did not complete within timeout")
	}
}

// TestTeeConcurrentSafety 测试并发安全性
func TestTeeConcurrentSafety(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	testData := strings.Repeat("concurrent test data ", 100)
	reader := strings.NewReader(testData)
	writer := newMockWriter()

	var processedCount int64
	var mu sync.Mutex

	handle := func(ctx context.Context, data []byte) error {
		mu.Lock()
		processedCount++
		mu.Unlock()
		// 模拟一些处理时间
		time.Sleep(time.Microsecond)
		return nil
	}

	tee := NewTee(ctx, logger, reader, writer, handle)
	defer tee.Close()

	err := tee.Stream()
	if err != nil {
		t.Fatalf("Stream() failed: %v", err)
	}

	// 等待所有数据处理完成
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	finalCount := processedCount
	mu.Unlock()

	if finalCount == 0 {
		t.Error("No data was processed")
	}

	t.Logf("Processed %d chunks of data", finalCount)
}

// TestBufferPoolEfficiency 测试缓冲区池的效率
func TestBufferPoolEfficiency(t *testing.T) {
	// 这个测试验证缓冲区池是否正常工作
	// 通过多次获取和归还缓冲区来测试

	var buffers []*[]byte

	// 获取多个缓冲区
	for i := 0; i < 10; i++ {
		bufPtr := bufferPool.Get().(*[]byte)
		buffers = append(buffers, bufPtr)

		// 验证缓冲区大小
		if len(*bufPtr) != 4096 {
			t.Errorf("Expected buffer size 4096, got %d", len(*bufPtr))
		}
	}

	// 归还所有缓冲区
	for _, bufPtr := range buffers {
		bufferPool.Put(bufPtr)
	}

	// 再次获取缓冲区，应该重用之前的缓冲区
	for i := 0; i < 5; i++ {
		bufPtr := bufferPool.Get().(*[]byte)
		if len(*bufPtr) != 4096 {
			t.Errorf("Expected reused buffer size 4096, got %d", len(*bufPtr))
		}
		bufferPool.Put(bufPtr)
	}
}

// BenchmarkTeeStream 基准测试
func BenchmarkTeeStream(b *testing.B) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))

	testData := strings.Repeat("benchmark test data ", 1000)

	handle := func(ctx context.Context, data []byte) error {
		return nil
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(testData)
		writer := io.Discard

		tee := NewTee(ctx, logger, reader, writer, handle)
		err := tee.Stream()
		if err != nil {
			b.Fatalf("Stream() failed: %v", err)
		}
		tee.Close()
	}
}

// BenchmarkBufferPool 缓冲区池基准测试
func BenchmarkBufferPool(b *testing.B) {
	b.Run("WithPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bufPtr := bufferPool.Get().(*[]byte)
			// 模拟使用缓冲区
			_ = *bufPtr
			bufferPool.Put(bufPtr)
		}
	})

	b.Run("WithoutPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := make([]byte, 4096)
			// 模拟使用缓冲区
			_ = buf
		}
	})
}

// TestTeeClose 测试关闭功能
func TestTeeClose(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	reader := strings.NewReader("test data")
	writer := newMockWriter()

	handle := func(ctx context.Context, data []byte) error {
		return nil
	}

	tee := NewTee(ctx, logger, reader, writer, handle)

	// 测试多次关闭不会 panic
	tee.Close()
	tee.Close()
	tee.Close()
}
