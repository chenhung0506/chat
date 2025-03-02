package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketConn struct {
	clients    map[string]*websocket.Conn
	services   map[string]*websocket.Conn
	clientsMu  sync.Mutex
	servicesMu sync.Mutex
}

var Conn = WebsocketConn{
	clients:  make(map[string]*websocket.Conn),
	services: make(map[string]*websocket.Conn),
}

func (s *WebsocketConn) AddClient(id string, conn *websocket.Conn) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	s.clients[id] = conn
}

func (s *WebsocketConn) GetClient(id string) (*websocket.Conn, bool) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	conn, ok := s.clients[id]
	return conn, ok
}

func (s *WebsocketConn) AddService(id string, conn *websocket.Conn) {
	s.servicesMu.Lock()
	defer s.servicesMu.Unlock()
	s.clients[id] = conn
}

func (s *WebsocketConn) GetService(id string) (*websocket.Conn, bool) {
	s.servicesMu.Lock()
	defer s.servicesMu.Unlock()
	conn, ok := s.clients[id]
	return conn, ok
}
