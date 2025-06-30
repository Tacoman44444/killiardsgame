import { Puck } from "./GameObjects"
import { JoinRoomData, ServerMessage } from "./socket-manager"
import { SocketEventManager } from "./socketevent-manager"
import { Board } from "./ui"
import { World } from "./World"


interface GameState {
    name:       string
    board:      Board
    socketEventBus: SocketEventManager
    processInput(input: any) : void;
    update() : void;
    render(ctx: CanvasRenderingContext2D) : void;
    Enter() : void;
}

export type GameInput = 
    | { type: "keydown", event: KeyboardEvent }
    | { type: "keyup", event: KeyboardEvent }
    | { type: "mousedown", event: MouseEvent }
    | { type: "mouseup", event: MouseEvent }
    | { type: "mousemove", event: MouseEvent };

class MainMenuState implements GameState {
    name:   string;
    board:  Board;
    socketEventBus: SocketEventManager

    constructor(socketEventBus: SocketEventManager) {
        this.name = "main-menu";
        this.board = new Board();
        this.socketEventBus = socketEventBus;
    }

    processInput(input: any) {
        this.board.processInput(input)
    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D) {
        this.board.render(ctx)
    }

    Enter() {}

    //handle socket event and send messages now ... 

    sendJoinRequest(id: string, code: string) {
        const data: JoinRoomData = { id: id, code: code }
        this.socketEventBus.emit("join-room", data)
    }

    sendCreateRoomRequest(id: string) {
        this.socketEventBus.emit("create-room", id)
    }
}

class RequestedForLobby implements GameState {
    name:   string;
    board:  Board;
    socketEventBus: SocketEventManager

    constructor(socketEventBus: SocketEventManager) {
        this.name = "requested-for-lobby";
        this.board = new Board();
        this.socketEventBus = socketEventBus;
    }

    processInput(input: any) {
        this.board.processInput(input)
    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D) {
        this.board.render(ctx)
    }

    Enter() {}

    //socket msges and send

    onInvalidCode(msg: ServerMessage) {
        //switch to MainMenuState
    }

    onRoomCreated(msg: ServerMessage) {
        //switch to InLobbyState
    }

    onRoomJoined(msg: ServerMessage) {
        //switch to InLobbyState
    }

}

class InLobby implements GameState {
    name:   string;
    board:  Board;
    socketEventBus: SocketEventManager

    constructor(socketEventBus: SocketEventManager) {
        this.name = "in-lobby";
        this.board = new Board();
        this.socketEventBus = socketEventBus;
    }

    processInput(input: any) {
        this.board.processInput(input)
    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D) {
        this.board.render(ctx)
    }

    Enter() {}

    //socket n send

    onGameStarted(msg: ServerMessage) {
        //switch to InGame state
    }

    sendGameStart() {
        this.socketEventBus.emit("start-game")
    }
}

class InGame implements GameState {
    name:   string;
    board:  Board;
    world: World | null;
    socketEventBus: SocketEventManager

    constructor(socketEventBus: SocketEventManager) {
        this.name = "in-lobby";
        this.board = new Board();
        this.world = null;
        this.socketEventBus = socketEventBus;
    }

    processInput(input: any) {
        this.board.processInput(input)
        this.world?.processInput(input)
    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D) {
        this.board.render(ctx)
    }

    Enter() {}

    initializeMap(msg: ServerMessage) {
        if (msg.type == "initialize-game") {
            let player = new Puck(msg.player.position_x, msg.player.position_y)
            let opps: Puck[] = []
            msg.other_players.forEach((opp) => opps.push(new Puck(opp.position_x, opp.position_y)))
            this.world = new World(msg.map.map, player, opps, this.socketEventBus)
        }
    }

    //handle them socket events and send them messages

    onGameOver(msg: ServerMessage) {}

}


class Game {
    canvas: HTMLCanvasElement;
    ctx: CanvasRenderingContext2D;
    currentState: GameState;
    socketEventBus: SocketEventManager;
    states: {
        mainMenuState: GameState;
        requestedForLobby: GameState;
        inLobby: GameState;
        inGame: GameState;
    }

    constructor(socketEventBus: SocketEventManager) {
        this.canvas = document.getElementById('gameCanvas') as HTMLCanvasElement;
        this.ctx = this.canvas.getContext('2d')!;
        this.socketEventBus = socketEventBus;
        this.states = {
            mainMenuState: new MainMenuState(socketEventBus),
            requestedForLobby: new RequestedForLobby(socketEventBus),
            inLobby: new InLobby(socketEventBus),
            inGame: new InGame(socketEventBus),
        }
        this.currentState = this.states.mainMenuState;
    }

    Play() {
    const loop = () => {
        this.currentState.update();
        this.currentState.render(this.ctx);
        requestAnimationFrame(loop);
        }
    requestAnimationFrame(loop);  
    }

    setupInputListeners() {
    window.addEventListener("keydown", (e) => {
        this.currentState.processInput({ type: "keydown", event: e });
    });

    window.addEventListener("keyup", (e) => {
        this.currentState.processInput({ type: "keyup", event: e });
    });

    this.canvas.addEventListener("mousedown", (e) => {
        this.currentState.processInput({ type: "mousedown", event: e });
    });

    this.canvas.addEventListener("mouseup", (e) => {
        this.currentState.processInput({ type: "mouseup", event: e });
    });

    this.canvas.addEventListener("mousemove", (e) => {
        this.currentState.processInput({ type: "mousemove", event: e });
    });
}

}
