package logit

import (
	"io"
	"sync"
	"sync/atomic"
)

type core struct {

	// level is the position of log.
	// In this version of logit, there are five levels:
	//
	//  DebugLevel, InfoLevel, WarnLevel, ErrorLevel, OffLevel.
	//
	// Higher level has higher visibility which means
	// one debug log will not be logged in one Logger set to InfoLevel.
	// That's we called level-based logger.
	//
	// In particular, OffLevel is the highest level, so if you set one
	// logger to OffLevel, it will shut up and log nothing.
	level *atomic.Value

	// needCaller is a flag to check if logs should contain caller's info or not.
	// This feature is useful but expensive in performance, so set to false if you don't need it.
	needCaller *atomic.Value

	// encoders are used to encode a log to bytes.
	// Every level has own encoder.
	encoders *encoders

	// writers are used to output an encoded log.
	// Every level has own writer.
	writers *writers
}

// newCore returns a new core for use.
func newCore(encoder Encoder, writer io.Writer) *core {
	return &core{
		level:      &atomic.Value{},
		needCaller: &atomic.Value{},
		encoders:   newEncoders(encoder),
		writers:    newWriters(writer),
	}
}

// Level returns the level of this core.
func (c *core) Level() Level {
	return c.level.Load().(Level)
}

// SetLevel will change the level to newLevel.
func (c *core) SetLevel(newLevel Level) {
	c.level.Store(newLevel)
}

// NeedCaller return needCaller of this core.
func (c *core) NeedCaller() bool {
	return c.needCaller.Load().(bool)
}

// SetNeedCaller sets needCaller to new one flag.
// If true, then every log will contain file name and line number.
// However, you should know that this is expensive in time.
// So be sure you really need it or keep it false.
func (c *core) SetNeedCaller(needCaller bool) {
	c.needCaller.Store(needCaller)
}

func (c *core) Encoders() *encoders {
	return c.encoders
}

func (c *core) Writers() *writers {
	return c.writers
}

// ============================================================================

type encoders struct {
	debugEncoder Encoder
	infoEncoder  Encoder
	warnEncoder  Encoder
	errorEncoder Encoder

	// lock is for safe concurrency.
	lock *sync.RWMutex
}

func newEncoders(encoder Encoder) *encoders {
	result := &encoders{
		debugEncoder: encoder,
		infoEncoder:  encoder,
		warnEncoder:  encoder,
		errorEncoder: encoder,
		lock:         &sync.RWMutex{},
	}
	return result
}

func (es *encoders) of(level Level) Encoder {

	es.lock.RLock()
	defer es.lock.RUnlock()

	if level == DebugLevel {
		return es.debugEncoder
	}

	if level == InfoLevel {
		return es.infoEncoder
	}

	if level == WarnLevel {
		return es.warnEncoder
	}

	return es.errorEncoder
}

// SetEncoder sets encoder to new one.
// This encoder will apply to all levels.
func (es *encoders) SetEncoder(encoder Encoder) {
	es.lock.Lock()
	defer es.lock.Unlock()
	es.debugEncoder = encoder
	es.infoEncoder = encoder
	es.warnEncoder = encoder
	es.errorEncoder = encoder
}

// SetDebugEncoder sets encoder of debug to new one.
func (es *encoders) SetDebugEncoder(encoder Encoder) {
	es.lock.Lock()
	defer es.lock.Unlock()
	es.debugEncoder = encoder
}

// SetInfoEncoder sets encoder of info to new one.
func (es *encoders) SetInfoEncoder(encoder Encoder) {
	es.lock.Lock()
	defer es.lock.Unlock()
	es.infoEncoder = encoder
}

// SetWarnEncoder sets encoder of warn to new one.
func (es *encoders) SetWarnEncoder(encoder Encoder) {
	es.lock.Lock()
	defer es.lock.Unlock()
	es.warnEncoder = encoder
}

// SetErrorEncoder sets encoder of error to new one.
func (es *encoders) SetErrorEncoder(encoder Encoder) {
	es.lock.Lock()
	defer es.lock.Unlock()
	es.errorEncoder = encoder
}

// ============================================================================

type writers struct {
	debugWriter io.Writer
	infoWriter  io.Writer
	warnWriter  io.Writer
	errorWriter io.Writer

	// lock is for safe concurrency.
	lock *sync.RWMutex
}

func newWriters(writer io.Writer) *writers {
	result := &writers{
		debugWriter: writer,
		infoWriter:  writer,
		warnWriter:  writer,
		errorWriter: writer,
		lock:        &sync.RWMutex{},
	}
	return result
}

func (ws *writers) of(level Level) io.Writer {

	ws.lock.RLock()
	defer ws.lock.RUnlock()

	if level == DebugLevel {
		return ws.debugWriter
	}

	if level == InfoLevel {
		return ws.infoWriter
	}

	if level == WarnLevel {
		return ws.warnWriter
	}

	return ws.errorWriter
}

// SetWriter sets writer to new one.
// This writer will apply to all levels.
func (ws *writers) SetWriter(writer io.Writer) {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	ws.debugWriter = writer
	ws.infoWriter = writer
	ws.warnWriter = writer
	ws.errorWriter = writer
}

// SetDebugWriter sets writer of debug to new one.
func (ws *writers) SetDebugWriter(writer io.Writer) {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	ws.debugWriter = writer
}

// SetInfoWriter sets writer of info to new one.
func (ws *writers) SetInfoWriter(writer io.Writer) {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	ws.infoWriter = writer
}

// SetWarnWriter sets writer of warn to new one.
func (ws *writers) SetWarnWriter(writer io.Writer) {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	ws.warnWriter = writer
}

// SetErrorWriter sets writer of error to new one.
func (ws *writers) SetErrorWriter(writer io.Writer) {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	ws.errorWriter = writer
}
