package main

type HubMessageType int

const (
	SendPlayerToLobby HubMessageType = iota
)

type HubMessage struct {
	msgType HubMessageType
}
