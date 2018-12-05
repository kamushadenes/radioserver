package StateModels

import (
	"github.com/racerxdl/radioserver/frontends"
	"github.com/racerxdl/radioserver/protocol"
	"sync"
)

type ServerState struct {
	DeviceInfo    protocol.DeviceInfo
	clients       []*ClientState
	clientListMtx sync.Mutex
	Frontend      frontends.Frontend
}

func CreateServerState() *ServerState {
	return &ServerState{
		clientListMtx: sync.Mutex{},
		clients:       make([]*ClientState, 0),
	}
}

func (s *ServerState) indexOfClient(state *ClientState) int {
	for k, v := range s.clients {
		if v.UUID == state.UUID {
			return k
		}
	}

	return -1
}

func (s *ServerState) PushClient(state *ClientState) {
	s.clientListMtx.Lock()
	defer s.clientListMtx.Unlock()

	s.clients = append(s.clients, state)
}

func (s *ServerState) RemoveClient(state *ClientState) {
	s.clientListMtx.Lock()
	defer s.clientListMtx.Unlock()
	idx := s.indexOfClient(state)
	if idx != -1 {
		s.clients = append(s.clients[:idx], s.clients[idx+1:]...)
	}
}

func (s *ServerState) SendSync() bool {
	s.clientListMtx.Lock()
	defer s.clientListMtx.Unlock()

	for i := 0; i < len(s.clients); i++ {
		go s.clients[i].SendSync()
	}

	return true
}

func (s *ServerState) PushSamples(samples []complex64) {

}
