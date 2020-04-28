package main

import (
	"fmt"
	"math/rand"
)

type message struct {
	data []byte
	room string
	name string
}

type subscription struct {
	conn *connection
	room string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var h = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

func removeValue(s []string, value string) []string {
	var index = -1
	for i, element := range listRoom {
		if element == value {
			index = i
			break
		}
	}
	if index == -1 {
		return s
	} else {
		s[index] = s[len(s)-1]
		return s[:len(s)-1]
	}

}

func (h *hub) run() {
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room]

			listRoom = append(listRoom, s.room)
			fmt.Println("room join:" + s.room)

			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
			}
			h.rooms[s.room][s.conn] = true
		case s := <-h.unregister:
			connections := h.rooms[s.room]
			fmt.Println("room leave:" + s.room)
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
						listRoom = removeValue(listRoom, s.room)
					}
				}
			}
		case m := <-h.broadcast:
			connections := h.rooms[m.room]
			for c := range connections {
				if c.id != m.name {
					select {
					case c.send <- m.data:
					default:
						close(c.send)
						delete(connections, c)
						if len(connections) == 0 {
							delete(h.rooms, m.room)
						}
					}
				}
			}
		}
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
