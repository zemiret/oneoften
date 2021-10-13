package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type ClientState struct {
	Lives int `json:"lives"`
	SeqNumber int `json:"seqNumber"`
}

type MessageType string

const (
	MessageDecreaseLive MessageType = "DECREASE_LIVE"
	MessagePlayerState MessageType = "PLAYER_STATE"
	MessageBuzzer MessageType = "MESSAGE_BUZZER"
)

type Message struct {
	MessageType MessageType `json:"messageType"`
	Payload     interface{} `json:"payload"`
}

type BuzzerPayload struct {
	Timestamp time.Time `json:"timestamp"`
}

type BuzzerResponse struct {
	SeqNumber int `json:"seqNumber"`
}

type Hub struct {
	sync.RWMutex

	log *log.Logger

	inbound chan InboundMessage
	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	state map[*Client]*ClientState
	lastTimestamp time.Time
}

func NewHub() *Hub {
	return &Hub{
		log:        log.New(os.Stdout, "[Hub]", log.Flags()),
		inbound:    make(chan InboundMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		state:      make(map[*Client]*ClientState),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.Lock()
			h.state[client] = &ClientState{
				Lives: 3,
				SeqNumber: client.seqNumber,
			}
			h.Unlock()

			if err := client.SendMessage(
				Message{
					MessageType: MessagePlayerState,
					Payload: h.state[client],
				}); err != nil {
				h.log.Printf("Err client.SendMessage %w", err)
			}

			//var roomForClient *Room
			//if client.gameType == GameTypeDuel {
			//	h.log.Printf("Adding %v to duel room\n", client.nickname)
			//	roomForClient = h.duelWaitingRoom
			//} else if client.gameType == GameTypeRoyale {
			//	h.log.Printf("Adding %v to royale room\n", client.nickname)
			//	roomForClient = h.royaleWaitingRoom
			//} else {
			//	h.log.Printf("Invalid gameType %v for %v\n", client.gameType, client.nickname)
			//	continue
			//}
			//
			//roomForClient.AddClient(client)
			//h.clients[client] = roomForClient
			//
			//if roomForClient.Full() {
			//	go roomForClient.Start()
			//
			//	if client.gameType == GameTypeDuel {
			//		h.duelWaitingRoom = newRoom(h.roomClosing, h.clientClosing, 2)
			//	} else if client.gameType == GameTypeRoyale {
			//		h.royaleWaitingRoom = newRoom(h.roomClosing, h.clientClosing, battleRoyaleRoomCapacity)
			//	}
			//}
		//case closingRoom := <-h.roomClosing:
		//	for client, room := range h.clients {
		//		if room != closingRoom {
		//			continue
		//		}
		//		h.closeClient(client)
		//	}
		//case closingClient := <-h.clientClosing:
		//h.closeClient(closingClient)
		case client := <-h.unregister:
			h.Lock()
			delete(h.state, client)
			h.Unlock()
		case inboundMessage := <-h.inbound:
			h.log.Println("Inbound message: ", string(inboundMessage.Message))

			client := inboundMessage.Client

			//if clientRoom, ok := h.clients[client]; ok {
			var messageType struct {
				MessageType MessageType `json:"messageType"`
			}
			if err := json.Unmarshal(inboundMessage.Message, &messageType); err != nil {
				h.log.Println("Inbound message malformed", err)
				continue
			}

			switch messageType.MessageType {
			case MessageDecreaseLive:
				h.Lock()
				if h.state[client].Lives > 0 {
					h.state[client].Lives -= 1
				}

				if err := client.SendMessage(
					Message{
						MessageType: MessagePlayerState,
						Payload: h.state[client],
					}); err != nil {
					h.log.Printf("Err client.SendMessage %w", err)
				}
				h.Unlock()
			case MessageBuzzer:
				var buzzerPayload BuzzerPayload
				if err := json.Unmarshal(inboundMessage.Message, &buzzerPayload); err != nil {
					h.log.Println("Invalid message", err)
					continue
				}

				h.Lock()
				if h.lastTimestamp.Before(time.Now().Add(-5 * time.Second)) ||
					buzzerPayload.Timestamp.Before(h.lastTimestamp) {

					h.lastTimestamp = buzzerPayload.Timestamp

					for cForResp, _ := range h.state {
						if err := cForResp.SendMessage(
							Message{
								MessageType: MessageBuzzer,
								Payload: BuzzerResponse{SeqNumber: client.seqNumber},
							}); err != nil {

							h.log.Printf("Err client.SendMessage %w", err)
						}
					}
				}
				h.Unlock()
			}

			//	switch messageType.MessageType {
			//	case MessageTypeBuyUnit:
			//		var buyUnitMessage BuyUnitMessage
			//		if err := json.Unmarshal(inboundMessage.Message, &buyUnitMessage); err != nil {
			//			h.log.Println("Invalid message", err)
			//			continue
			//		}
			//		buyUnitMessage.Client = inboundMessage.Client
			//		clientRoom.BuyUnitChannel <- buyUnitMessage
			//	case MessageTypeSellUnit:
			//		var sellUnitMessage SellUnitMessage
			//		if err := json.Unmarshal(inboundMessage.Message, &sellUnitMessage); err != nil {
			//			h.log.Println("Invalid message", err)
			//			continue
			//		}
			//		sellUnitMessage.Client = inboundMessage.Client
			//		clientRoom.SellUnitChannel <- sellUnitMessage
			//	case MessageTypePlaceUnit:
			//		var placeUnitMessage PlaceUnitMessage
			//		if err := json.Unmarshal(inboundMessage.Message, &placeUnitMessage); err != nil {
			//			h.log.Println("Invalid message", err)
			//			continue
			//		}
			//		placeUnitMessage.Client = inboundMessage.Client
			//		clientRoom.PlaceUnitChannel <- placeUnitMessage
			//	case MessageTypeUnplaceUnit:
			//		var unplaceUnitMessage UnplaceUnitMessage
			//		if err := json.Unmarshal(inboundMessage.Message, &unplaceUnitMessage); err != nil {
			//			h.log.Println("Invalid message", err)
			//			continue
			//		}
			//		unplaceUnitMessage.Client = inboundMessage.Client
			//		clientRoom.UnplaceUnitChannel <- unplaceUnitMessage
			//	default:
			//		h.log.Println("Unknown message type")
			//	}
			//}
		}
	}
}
