package main

import (
	"fmt"
	"time"

	"github.com/Tacoman44444/killiardsgame/server/tools"
	"github.com/gorilla/websocket"
)

type LobbyMessageType int

const (
	LobbySendEntityUpdate LobbyMessageType = iota
	LobbySendWallUpdate
	LobbySendMapUpdate
	LobbySendTurnStart
	LobbySendTurnTimeout
	LobbyBroadcastMove
	LobbySendGameStart
	LobbyClose // :<
)

type LobbyMessage struct {
	msgType    LobbyMessageType
	lobbyCode  string
	player     PlayerIdentity
	allPlayers []PlayerIdentity
	walls      []WallState
	mapState   tools.MapState
}

type LobbyState interface {
	Enter(lobby *Lobby)
	HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby)
	HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby)
	Exit(lobby *Lobby)
}

type Lobby struct {
	code      string
	Inbound   chan PlayerMessage
	readHub   chan HubMessage
	state     LobbyState
	gameState *GameState
	players   map[*websocket.Conn]*Player
	turnCycle []*websocket.Conn
	turn      int
	owner     *Player
}

func NewLobby(code string, owner *Player) *Lobby {
	lb := Lobby{
		code:      code,
		Inbound:   make(chan PlayerMessage),
		readHub:   make(chan HubMessage),
		state:     &LobbyWaitingForPlayers{},
		gameState: nil,
		players:   make(map[*websocket.Conn]*Player),
		turnCycle: make([]*websocket.Conn, 0, 10),
		turn:      0,
		owner:     owner,
	}
	lb.turnCycle = append(lb.turnCycle, owner.conn)
	lb.players[owner.conn] = owner

	return &lb
}

func (l *Lobby) SetState(newState LobbyState) {
	l.state.Exit(l)
	l.state = newState
	l.state.Enter(l)
}

func (l *Lobby) SwitchTurn() {
	if l.turn == len(l.turnCycle)-1 {
		l.turn = 0
	} else {
		l.turn++
	}
}

type LobbyWaitingForPlayers struct{}

func (l *LobbyWaitingForPlayers) Enter(lobby *Lobby) {

}

func (l *LobbyWaitingForPlayers) HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby) {
	if !channelOpen {
		//handle player channel not being open
	}

	switch pm.msgType {
	case PlayerStartGame:
		if pm.sender == lobby.owner.conn {
			//start game and send all clients a message that the game has begun
			for key, _ := range lobby.players {
				lobby.turnCycle = append(lobby.turnCycle, key)
			}
			lobby.SetState(&LobbyWaitingForTurn{})
			players := make([]string, 10)
			for _, value := range lobby.players {
				players = append(players, value.id)
			}
			lobby.gameState = GetNewGame(players)

			for _, value := range lobby.players { //sending the message to all players
				msg := LobbyMessage{
					msgType:    LobbySendGameStart,
					player:     *lobby.gameState.players[value.id],
					allPlayers: PlayerMapToSlice(lobby.gameState.players),
					walls:      WallStateRefToWallState(lobby.gameState.walls),
					mapState:   *lobby.gameState.mapState,
				}
				value.readLobby <- msg
			}
		} else if pm.sender != lobby.owner.conn {
			fmt.Println("only the party owner can start the match")
			return
		}
	}
}

func (l *LobbyWaitingForPlayers) HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby) {
	if !channelOpen {
		//handle player channel not being open
	}

	switch hm.msgType {
	case HubSendPlayerToLobby:
		lobby.players[hm.player.conn] = hm.player
	}
}

func (l *LobbyWaitingForPlayers) Exit(lobby *Lobby) {

}

type LobbyWaitingForTurn struct {
	timer  *time.Timer
	cancel bool
}

func (l *LobbyWaitingForTurn) Enter(lobby *Lobby) {
	l.timer = time.NewTimer(lobby.gameState.turnTimer)

	msg := LobbyMessage{
		msgType:    LobbySendTurnStart,
		player:     PlayerIdentity{},
		allPlayers: nil,
		walls:      nil,
		mapState:   tools.MapState{},
	}
	lobby.players[lobby.turnCycle[lobby.turn]].readLobby <- msg

	go func() {
		<-l.timer.C
		if l.cancel {
			return
		}
		msg := LobbyMessage{
			msgType:    LobbySendTurnTimeout,
			player:     PlayerIdentity{},
			allPlayers: PlayerMapToSlice(lobby.gameState.players),
			walls:      WallStateRefToWallState(lobby.gameState.walls),
			mapState:   *lobby.gameState.mapState,
		}
		lobby.players[lobby.turnCycle[lobby.turn]].readLobby <- msg
		lobby.SwitchTurn()
		//change state
		lobby.SetState(&LobbyWaitingForTurn{})

	}()
}

func (l *LobbyWaitingForTurn) HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby) {
	if !channelOpen {
		//handle channel being closed
	}

	switch pm.msgType {
	case PlayerStartGame:
		fmt.Println("game already in progress")
	case PlayerSendAction:
		if pm.sender == lobby.turnCycle[lobby.turn] {
			for _, value := range lobby.players {
				if pm.player.id != value.id {
					msg := LobbyMessage{
						msgType:    LobbyBroadcastMove,
						player:     *lobby.gameState.players[value.id],
						allPlayers: PlayerMapToSlice(lobby.gameState.players),
						walls:      WallStateRefToWallState(lobby.gameState.walls),
						mapState:   *lobby.gameState.mapState,
					}
					value.readLobby <- msg
				}
			}
			tools.PhysicsResolver(lobby.gameState.players[pm.senderID].circle, PlayerMapToCircles(lobby.gameState.players), GetWallRectRefs(lobby.gameState.walls), pm.msg.Action.data)
			for _, value := range lobby.players { //sending the message to all players
				msg := LobbyMessage{
					msgType:    LobbySendEntityUpdate,
					player:     *lobby.gameState.players[value.id],
					allPlayers: PlayerMapToSlice(lobby.gameState.players),
					walls:      WallStateRefToWallState(lobby.gameState.walls),
					mapState:   *lobby.gameState.mapState,
				}
				value.readLobby <- msg
			}

			lobby.SetState(&LobbyBetweenTurns{})
		}
	case PlayerSendWall:
		if pm.sender == lobby.turnCycle[lobby.turn] {
			newWall := pm.msg.Wall
			lobby.gameState.walls = append(lobby.gameState.walls, &newWall)
			for _, value := range lobby.players { //sending the message to all players
				msg := LobbyMessage{
					msgType:    LobbySendWallUpdate,
					player:     *lobby.gameState.players[value.id],
					allPlayers: PlayerMapToSlice(lobby.gameState.players),
					walls:      WallStateRefToWallState(lobby.gameState.walls),
					mapState:   *lobby.gameState.mapState,
				}
				value.readLobby <- msg
			}
		}
	}
}

func (l *LobbyWaitingForTurn) HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby) {

}

func (l *LobbyWaitingForTurn) Exit(lobby *Lobby) {
	l.cancel = true
}

type LobbyBetweenTurns struct {
	timer      *time.Timer //this timer is for when clients dont send "simulation-done" in time
	cancel     bool
	simRunning []*websocket.Conn
}

func (l *LobbyBetweenTurns) Enter(lobby *Lobby) {
	l.timer = time.NewTimer(time.Duration(CLIENT_AFFIRMATON_TIMEOUT_IN_SECONDS) * time.Second)
	l.simRunning = lobby.turnCycle

	go func() {
		<-l.timer.C
		if l.cancel {
			return
		}
		lobby.SwitchTurn()
		msg := LobbyMessage{
			msgType:    LobbySendTurnStart,
			player:     *lobby.gameState.players[lobby.players[lobby.turnCycle[lobby.turn]].id],
			allPlayers: nil,
			walls:      WallStateRefToWallState(lobby.gameState.walls),
			mapState:   *lobby.gameState.mapState,
		}
		for _, value := range lobby.players {
			value.readLobby <- msg
		}
		lobby.SetState(&LobbyWaitingForTurn{})

	}()
}

func (l *LobbyBetweenTurns) HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby) {
	if !channelOpen {
		//handle channel not being open
	}

	switch pm.msgType {
	case PlayerSimulationDone:
		for i := range l.simRunning {
			if l.simRunning[i] == pm.sender {
				l.simRunning = append(l.simRunning[:i], l.simRunning[i+1:]...)
			}
		}
		//HANDLE THE CASE FOR WHEN THE PLAYER DISCONNECTS
	}

	if len(l.simRunning) == 0 {
		lobby.SwitchTurn()
		msg := LobbyMessage{
			msgType:    LobbySendTurnStart,
			player:     *lobby.gameState.players[lobby.players[lobby.turnCycle[lobby.turn]].id],
			allPlayers: nil,
			walls:      WallStateRefToWallState(lobby.gameState.walls),
			mapState:   *lobby.gameState.mapState,
		}
		for _, value := range lobby.players {
			value.readLobby <- msg
		}

		lobby.SetState(&LobbyWaitingForTurn{})

	}
}

func (l *LobbyBetweenTurns) HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby) {

}

func (l *LobbyBetweenTurns) Exit(lobby *Lobby) {
	l.cancel = true
}
