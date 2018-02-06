// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package quic

import (
	"fmt"
	"sync"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/qerr"
)

type outgoingUniStreamsMap struct {
	mutex sync.RWMutex

	streams map[protocol.StreamID]sendStreamI

	nextStream protocol.StreamID
	newStream  func(protocol.StreamID) sendStreamI

	closeErr error
}

func newOutgoingUniStreamsMap(nextStream protocol.StreamID, newStream func(protocol.StreamID) sendStreamI) *outgoingUniStreamsMap {
	return &outgoingUniStreamsMap{
		streams:    make(map[protocol.StreamID]sendStreamI),
		nextStream: nextStream,
		newStream:  newStream,
	}
}

func (m *outgoingUniStreamsMap) OpenStream() (sendStreamI, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.closeErr != nil {
		return nil, m.closeErr
	}
	s := m.newStream(m.nextStream)
	m.streams[m.nextStream] = s
	m.nextStream += 4
	return s, nil
}

func (m *outgoingUniStreamsMap) GetStream(id protocol.StreamID) (sendStreamI, error) {
	if id >= m.nextStream {
		return nil, qerr.Error(qerr.InvalidStreamID, fmt.Sprintf("peer attempted to open stream %d", id))
	}
	m.mutex.RLock()
	s := m.streams[id]
	m.mutex.RUnlock()
	return s, nil
}

func (m *outgoingUniStreamsMap) DeleteStream(id protocol.StreamID) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.streams[id]; !ok {
		return fmt.Errorf("Tried to delete unknown stream %d", id)
	}
	delete(m.streams, id)
	return nil
}

func (m *outgoingUniStreamsMap) CloseWithError(err error) {
	m.mutex.Lock()
	m.closeErr = err
	m.mutex.Unlock()
}
