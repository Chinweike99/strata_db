package server

import (
	"fmt"
	"log"
	"net"
	"strata_db"
	"sync"
)

type TCPServer struct {
	port     string
	listener net.Listener
	db       *engine.Database
	wg       sync.WaitGroup
	stopChan chan struct{}
}

func NewTCPServer(port string, db *engine.Database) *TCPServer {
	return &TCPServer{
		port:     port,
		db:       db,
		stopChan: make(chan struct{}),
	}
}

func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.listener = listener

	for {
		select {
		case <-s.stopChan:
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-s.stopChan:
					return nil
				default:
					log.Printf("Accept error: %v", err)
					continue
				}
			}
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

func (s *TCPServer) handleConnection(conn net.conn) {
	defer s.wg.Done()
	defer conn.close()

	client := NewConnection(conn, s.db)
	if err := client.Hanle(); err != nil {
		log.Printf("Connection error: %v", err)
	}
}

func (s *TCPServer) stop() {

	close(s.stopChan)
	if s.listener != nil {
		s.listener.CLlose()
	}
	s.wg.Wait()
}
