package lumberjack

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

// BufferedFile is bufferChan writer than can be reopned
type BufferedFile struct {
	mutex          sync.Mutex
	quitChan       chan bool
	done           bool
	async          bool
	bufferChan     chan *bufferedMsg
	quitBufferChan chan bool
	OrigRLog       *os.File
	BufWriter      *bufio.Writer
}

type bufferedMsg struct {
	data []byte
	quit bool
}

var (
	// defaultBufferSize exists so it can be mocked out by tests
	defaultBufferSize = 4 * 1024 * 1024
	// defaultFlushInterval exists so it can be mocked out by tests
	defaultFlushInterval = 1 * time.Second

	defaultChannelSize = 32 * 1024
)

// NewBufferedFile opens a buffered file that is periodically flushed.
func NewBufferedFile(rl *os.File, bufferSize, channelSize int, asyncWrite bool) *BufferedFile {
	if bufferSize <= 0 {
		bufferSize = defaultBufferSize
	}
	if channelSize <= 0 {
		channelSize = defaultChannelSize
	}
	return NewBufferedFileWithDetails(rl, bufferSize, defaultFlushInterval, channelSize, asyncWrite)
}

// NewBufferedFileWithDetails opens a buffered file with the
// given buffer size that is periodically flushed on the given interval
// or given channel size that is asynchronous write logs when asyncWrite is on
func NewBufferedFileWithDetails(rl *os.File, bufferSize int, flush time.Duration, channelSize int, asyncWrite bool) *BufferedFile {
	brl := BufferedFile{
		quitChan:  make(chan bool, 1),
		async:     asyncWrite,
		OrigRLog:  rl,
		BufWriter: bufio.NewWriterSize(rl, bufferSize),
	}
	if brl.async {
		brl.bufferChan = make(chan *bufferedMsg, channelSize)
		brl.quitBufferChan = make(chan bool, 0)
		go brl.writeDaemon()
	}
	go brl.flushDaemon(flush)
	return &brl
}

// flushDaemon periodically flushes the log file buffers
func (brl *BufferedFile) flushDaemon(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-brl.quitChan:
			ticker.Stop()
			return
		case <-ticker.C:
			brl.Flush()
		}
	}
}

// writeDaemon receive log data from the channel periodically and write synchronously
func (brl *BufferedFile) writeDaemon() {
	for {
		buf := <-brl.bufferChan
		if buf.quit {
			brl.quitBufferChan <- true
			break
		}
		brl.syncWrite(buf.data)
	}
}

// Flush flushes the bufferChan
func (brl *BufferedFile) Flush() {
	brl.mutex.Lock()

	if brl.done {
		brl.mutex.Unlock()
		return
	}

	brl.BufWriter.Flush()
	brl.OrigRLog.Sync()
	brl.mutex.Unlock()
}

// Write implements io.WriteCloser
func (brl *BufferedFile) Write(p []byte) (int, error) {
	if brl.async {
		return brl.asyncWrite(p)
	}
	return brl.syncWrite(p)
}

// asyncWrite write logs asynchronously through channel
func (brl *BufferedFile) asyncWrite(p []byte) (int, error) {
	b := make([]byte, len(p))
	copy(b, p)
	select {
	case brl.bufferChan <- &bufferedMsg{
		data: b,
		quit: false,
	}:
	default:
		return 0, fmt.Errorf("can't async write, chan size is full")
	}
	return len(b), nil
}

// syncWrite write logs to bufio synchronously
func (brl *BufferedFile) syncWrite(p []byte) (int, error) {
	brl.mutex.Lock()
	n, err := brl.BufWriter.Write(p)

	// means flush happenede in the middle of the line
	// and we need to flush the rest of our string at this point
	if brl.BufWriter.Buffered() < len(p) {
		brl.BufWriter.Flush()
	}

	brl.mutex.Unlock()
	return n, err
}

// Close flushes the internal bufferChan and closes the destination file
func (brl *BufferedFile) Close() error {
	brl.quitChan <- true
	if brl.async {
		brl.bufferChan <- &bufferedMsg{quit: true}
		<-brl.quitBufferChan
		close(brl.bufferChan)
	}
	brl.mutex.Lock()

	brl.done = true
	brl.BufWriter.Flush()

	err := brl.OrigRLog.Close()
	if err != nil {
		return err
	}

	brl.mutex.Unlock()
	return nil
}

func (brl *BufferedFile) Stat() (os.FileInfo, error) {
	return brl.OrigRLog.Stat()
}
