package main

import (
	"fmt"

	"github.com/Tacoman44444/killiardsgame/server/tools"
	"github.com/gorilla/websocket"
)

type LobbyMessageType int

const (
	LobbySendEntityUpdate LobbyMessageType = iota
	LobbySendEliminations
	LobbySendWallUpdate
	LobbySendMapUpdate
	LobbySendTurnStart
	LobbySendTurnTimeout
	LobbyBroadcastMove
	LobbySendGameStart
	LobbySendGameOver
	LobbyClose // :<
)

type LobbyMessage struct {
	msgType           LobbyMessageType
	lobbyCode         string
	player            PlayerIdentity
	allPlayers        []PlayerIdentity
	eliminatedPlayers []PlayerIdentity
	walls             []WallState
	currentMap        tools.MapState
	nextMap           tools.MapState
	action            PlayerAction
	result            string
	winnerName        string
}

type TurnQueue struct {
	data       []*Player
	currentIdx int
}

func NewTurnQueue() *TurnQueue {
	return &TurnQueue{
		data:       []*Player{},
		currentIdx: 0,
	}
}

func (q *TurnQueue) Add(p *Player) {
	q.data = append(q.data, p)
}

func (q *TurnQueue) RemoveByID(id string) bool {
	if len(q.data) == 0 {
		fmt.Println("queue is empty")
		return false
	}

	index := -1
	for i, p := range q.data {
		if p.id == id {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("player not found")
		return false
	}

	// Remove the player from the slice
	q.data = append(q.data[:index], q.data[index+1:]...)

	// Adjust currentIdx if needed
	if index < q.currentIdx {
		q.currentIdx-- // shift left to account for removed item
	} else if index == q.currentIdx {
		if q.currentIdx == 0 {
			q.currentIdx = len(q.data) - 1
		} else {
			q.currentIdx--
		}
	}

	// Handle edge case: queue is now empty
	if len(q.data) == 0 {
		q.currentIdx = 0
	}

	return true
}
func (q *TurnQueue) Current() *Player {
	return q.data[q.currentIdx]
}

func (q *TurnQueue) Next() *Player {
	q.currentIdx = (q.currentIdx + 1) % len(q.data)
	return q.data[q.currentIdx]
}

func (q *TurnQueue) List() []*Player {
	return q.data
}

func (q *TurnQueue) Size() int {
	return len(q.data)
}

type Lobby struct {
	code              string
	Inbound           chan PlayerMessage
	readHub           chan HubMessage
	gameState         *GameState
	players           map[*websocket.Conn]*Player
	queue             *TurnQueue
	owner             *Player
	eliminated        []*Player
	currentState      LobbyState
	waitingforplayers LobbyWaitingForPlayers
	inturn            LobbyInTurn
	processingturn    LobbyProcessingTurn
	gameOver          LobbyGameOver
	simCount          int
}

func NewLobby(code string, owner *Player) *Lobby {
	lb := Lobby{
		code:      code,
		Inbound:   make(chan PlayerMessage),
		readHub:   make(chan HubMessage),
		gameState: nil,
		players:   make(map[*websocket.Conn]*Player),
		queue:     NewTurnQueue(),
		owner:     owner,
		simCount:  0,
	}
	lb.waitingforplayers = LobbyWaitingForPlayers{}
	lb.inturn = LobbyInTurn{}
	lb.processingturn = LobbyProcessingTurn{}
	lb.gameOver = LobbyGameOver{}
	lb.currentState = lb.waitingforplayers
	lb.currentState.Enter(&lb)
	lb.players[owner.conn] = owner
	lb.queue.Add(owner)

	return &lb
}

func (l *Lobby) Run() {
	fmt.Printf("[LOBBY] Lobby %s is now running\n", l.code)
	for {
		select {
		case pm, ok := <-l.Inbound:
			if !ok {
				fmt.Printf("[LOBBY] Lobby %s Inbound channel closed\n", l.code)
				return
			}
			l.currentState.HandlePlayerMessage(pm, ok, l)

		case hm, ok := <-l.readHub:
			if !ok {
				fmt.Printf("[LOBBY] Lobby %s readHub channel closed\n", l.code)
				return
			}
			l.currentState.HandleHubMessage(hm, ok, l)
		}
	}
}

func (l *Lobby) NextTurn() {
	l.queue.Next()
}

func (l *Lobby) SetState(state LobbyState) {
	l.currentState.Exit(l)
	l.currentState = state
	l.currentState.Enter(l)
}

func (l *Lobby) Eliminate() []PlayerIdentity {
	activePlayers := append([]*Player{}, l.queue.List()...)
	eliminatedThisRound := make([]PlayerIdentity, 0, 10)
	for i := range activePlayers {
		if tools.IsPlayerEliminated(l.gameState.mapState, l.gameState.players[activePlayers[i].id].circle.Center) {
			//DELTE THE PLAYAA
			if l.queue.RemoveByID(activePlayers[i].id) {
				l.eliminated = append(l.eliminated, activePlayers[i])
				eliminatedThisRound = append(eliminatedThisRound, *l.gameState.players[activePlayers[i].id])
			}
		}
	}
	return eliminatedThisRound
}

type LobbyState interface {
	Enter(lobby *Lobby)
	HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby)
	HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby)
	Exit(lobby *Lobby)
}

type LobbyWaitingForPlayers struct{}

func (l LobbyWaitingForPlayers) Enter(lobby *Lobby) {}
func (l LobbyWaitingForPlayers) HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby) {
	if !channelOpen {
	}

	switch pm.msgType {
	case PlayerStartGame:
		if pm.sender == lobby.owner.conn {
			playerIDs := make([]string, 0, 10)
			playerUsernames := make([]string, 0, 10)
			for _, value := range lobby.players {
				playerIDs = append(playerIDs, value.id)
				playerUsernames = append(playerUsernames, value.username)
			}
			lobby.gameState = GetNewGame(playerIDs, playerUsernames)

			for _, value := range lobby.players { //sending the message to all players
				msg := LobbyMessage{
					msgType:    LobbySendGameStart,
					player:     *lobby.gameState.players[value.id],
					allPlayers: PlayerMapToSlice(lobby.gameState.players),
					walls:      WallStateRefToWallState(lobby.gameState.walls),
					currentMap: *lobby.gameState.mapState,
					nextMap:    *lobby.gameState.nextMap,
				}
				value.readLobby <- msg
			}
			lobby.SetState(lobby.inturn)
		} else if pm.sender != lobby.owner.conn {
			fmt.Println("only the party owner can start the match")
			return
		}
	}
}
func (l LobbyWaitingForPlayers) HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby) {
	if !channelOpen {
	}

	switch hm.msgType {
	case HubSendPlayerToLobby:
		lobby.players[hm.player.conn] = hm.player
		lobby.queue.Add(hm.player)
	}
}
func (l LobbyWaitingForPlayers) Exit(lobby *Lobby) {}

type LobbyInTurn struct{}

func (l LobbyInTurn) Enter(lobby *Lobby) {
	lobby.NextTurn()
	msg := LobbyMessage{
		msgType:    LobbySendTurnStart,
		player:     PlayerIdentity{},
		allPlayers: nil,
		walls:      nil,
	}
	player := lobby.queue.Current()
	player.turnsPlayed++
	fmt.Println("sending a turn start message to: ", player.id)
	player.readLobby <- msg
}
func (l LobbyInTurn) HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby) {
	if !channelOpen {
	}
	switch pm.msgType {
	case PlayerSendAction:
		fmt.Println("PLAYA sent an action")
		if pm.sender == lobby.queue.Current().conn {
			for _, value := range lobby.players {
				if pm.player.id != value.id {
					fmt.Println("BROADCASTING MOVE")
					msg := LobbyMessage{
						msgType: LobbyBroadcastMove,
						player:  *lobby.gameState.players[pm.player.id],
						action:  pm.msg.Action,
					}
					value.readLobby <- msg
				}
			}
			tools.PhysicsResolver(lobby.gameState.players[pm.senderID].circle, PlayerMapToCircles(lobby.gameState.players), GetWallRectRefs(lobby.gameState.walls), PlayerActionToShotData(pm.msg.Action))
			for _, value := range lobby.players {
				fmt.Println("Sending entity update message to: ", value.id)
				//turn the active queue into a list, then get them playeridentities
				activePlayers := append([]*Player{}, lobby.queue.List()...)
				activePlayerIDs := make([]PlayerIdentity, 0, 10)
				for i := range activePlayers {
					activePlayerIDs = append(activePlayerIDs, *lobby.gameState.players[activePlayers[i].id])
				}
				msg := LobbyMessage{
					msgType:    LobbySendEntityUpdate,
					player:     *lobby.gameState.players[value.id],
					allPlayers: activePlayerIDs,
					walls:      WallStateRefToWallState(lobby.gameState.walls),
				}
				value.readLobby <- msg
			}

			lobby.SetState(lobby.processingturn)
		} else {
			fmt.Println("the guy who send the move doesnt seem to match with the guy who should be the one sendinrgirngrihafug")
			fmt.Println("sender ID: ", pm.sender)
			fmt.Println("active guy ID: ", lobby.queue.Current().conn)
		}
	case PlayerSendWall:
		if pm.sender == lobby.queue.Current().conn {
			newWall := pm.msg.Wall
			lobby.gameState.walls = append(lobby.gameState.walls, &newWall)
			for _, value := range lobby.players {
				msg := LobbyMessage{
					msgType:    LobbySendWallUpdate,
					player:     *lobby.gameState.players[value.id],
					allPlayers: PlayerMapToSlice(lobby.gameState.players),
					walls:      WallStateRefToWallState(lobby.gameState.walls),
				}
				value.readLobby <- msg
			}
		}
	}
}
func (l LobbyInTurn) HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby) {}
func (l LobbyInTurn) Exit(lobby *Lobby)                                              {}

type LobbyProcessingTurn struct{}

func (l LobbyProcessingTurn) Enter(lobby *Lobby) {
	lobby.simCount = 0
}
func (l LobbyProcessingTurn) HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby) {
	if !channelOpen {
	}
	switch pm.msgType {
	case PlayerSimulationDone:
		lobby.simCount++
		fmt.Println("simCount is: ", lobby.simCount, " and the no. of players are: ", lobby.queue.Size())
		if lobby.simCount >= lobby.queue.Size() {
			//eliminate the dead players...
			//first we shrink the map if 4turns/8turns have happened
			minTurns := lobby.queue.List()[0].turnsPlayed
			for _, p := range lobby.queue.List() {
				if p.turnsPlayed < minTurns {
					minTurns = p.turnsPlayed
				}
			}

			if minTurns == 1 {
				fmt.Println("sending map update blyat")
				nextMap := tools.ShrinkArena(lobby.gameState.nextMap)
				lobby.gameState.mapState = lobby.gameState.nextMap
				lobby.gameState.nextMap = nextMap
				msg := LobbyMessage{
					msgType:    LobbySendMapUpdate,
					currentMap: *lobby.gameState.mapState,
					nextMap:    *lobby.gameState.nextMap,
				}
				for _, value := range lobby.players {
					value.readLobby <- msg
				}
			}
			if minTurns == 4 {
				//apply the shrinkmap
				lobby.gameState.mapState = lobby.gameState.nextMap
				msg := LobbyMessage{
					msgType:    LobbySendMapUpdate,
					currentMap: *lobby.gameState.mapState,
					nextMap:    *lobby.gameState.nextMap,
				}
				for _, value := range lobby.players {
					value.readLobby <- msg
				}
			}
			eliminated := lobby.Eliminate()
			if len(eliminated) != 0 {
				//some1 dead
				msg := LobbyMessage{
					msgType:           LobbySendEliminations,
					eliminatedPlayers: eliminated,
				}
				for _, value := range lobby.players {
					value.readLobby <- msg
				}

			} else {
				fmt.Println("everybody lived this turn")
			}
			if lobby.queue.Size() == 0 {
				//we have a draw
				msg := LobbyMessage{
					msgType:    LobbySendGameOver,
					result:     "draw",
					winnerName: "",
				}
				for _, value := range lobby.players {
					value.readLobby <- msg
				}
				lobby.SetState(lobby.gameOver)

			} else if lobby.queue.Size() == 1 {
				//ladies and gentlemen we have a winner
				msg := LobbyMessage{
					msgType:    LobbySendGameOver,
					result:     "win",
					winnerName: lobby.queue.Current().id,
				}
				for _, value := range lobby.players {
					value.readLobby <- msg
				}
				lobby.SetState(lobby.gameOver)
			} else {
				lobby.SetState(lobby.inturn)
			}

		}
	}
}
func (l LobbyProcessingTurn) HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby) {}
func (l LobbyProcessingTurn) Exit(lobby *Lobby)                                              { lobby.simCount = 0 }

type LobbyGameOver struct{}

func (l LobbyGameOver) Enter(lobby *Lobby) {}
func (l LobbyGameOver) HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby) {
	//handle player returning to lobby
}
func (l LobbyGameOver) HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby) {}
func (l LobbyGameOver) Exit(lobby *Lobby)                                              {}
