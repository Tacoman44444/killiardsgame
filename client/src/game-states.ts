import { SpriteComponent } from "./Components"
import { Puck } from "./GameObjects"
import { JoinRoomData, ServerMessage } from "./socket-manager"
import { SocketEventManager } from "./socketevent-manager"
import { Board } from "./ui"
import { World } from "./World"


interface GameState {
    name:       string
    board:      Board
    game: Game
    initializeBoard() : void;
    processInput(input: any) : void;
    render(ctx: CanvasRenderingContext2D) : void;
}

export enum MouseButton {
    LEFT_CLICK,
    MIDDLE_CLICK,
    RIGHT_CLICK,
}

export type GameInput = 
    | { type: "keydown", event: KeyboardEvent }
    | { type: "keyup", event: KeyboardEvent }
    | { type: "mousedown", event: MouseEvent, buttonType: MouseButton, cameraX: number, cameraY: number }
    | { type: "mouseup", event: MouseEvent, buttonType: MouseButton, cameraX: number, cameraY: number }
    | { type: "mousemove", event: MouseEvent, cameraX: number, cameraY: number }


class MainMenuState implements GameState {
    name:   string;
    board:  Board;
    game: Game;

    constructor(game: Game) {
        this.name = "main-menu";
        this.board = new Board();
        this.game = game;
    }

    initializeBoard() {
        this.board.addInputTextBox("username", 600, 525, 400, 150, 10);
        this.board.addInputTextBox("code", 700, 200, 200, 75, 6);
        this.board.addButton("create", 20, 20, 68, 16, new SpriteComponent("createroom_button"), () => this.sendCreateRoomRequest(), true);
        this.board.addButton("join", 20, 700, 65, 16, new SpriteComponent("createroom_button"), () => this.sendJoinRequest(), true)
    }

    processInput(input: any) {
        this.board.processInput(input)
    }

    render(ctx: CanvasRenderingContext2D) {
        this.board.render(ctx)
    }

    //handle socket event and send messages now ... 

    sendJoinRequest() {
        let code = "";
        this.board.textBoxes.forEach((box) => {
            if (box.name === "code") {
                code = box.text;
            }
        });
        this.game.socketEventBus.emit("join-room", code);
        this.game.currentState = this.game.states.requestedForLobby;
        this.game.states.requestedForLobby;
    }

    sendCreateRoomRequest() {
        this.game.socketEventBus.emit("create-room")
        this.game.currentState = this.game.states.requestedForLobby;
    }
}

class RequestedForLobby implements GameState {
    name:   string;
    board:  Board;
    game: Game;
    code: string;   //rage coding, probably shouldnt do it this way

    constructor(game: Game) {
        this.name = "requested-for-lobby";
        this.board = new Board();
        this.game = game;
        this.code = "";
    }

    initializeBoard() {

    }

    processInput(input: any) {
        this.board.processInput(input)
    }

    render(ctx: CanvasRenderingContext2D) {
        this.board.render(ctx)
    }

    //socket msges and send

    onInvalidCode(msg: ServerMessage) {
        this.game.currentState = this.game.states.mainMenuState;
    }

    onRoomCreated(msg: ServerMessage) {
        this.game.currentState = this.game.states.inLobby;
        if (msg.type == "room-created") {
            this.game.code = msg.code;
        }
    }

    onRoomJoined(msg: ServerMessage) {
        this.game.currentState = this.game.states.inLobby;
        if (msg.type == "room-joined") {
            this.game.code = msg.code;
        }
    }
}

class InLobby implements GameState {
    name:   string;
    board:  Board;
    game: Game;

    constructor(game: Game) {
        this.name = "in-lobby";
        this.board = new Board();
        this.game = game;
    }

    initializeBoard() {
        this.board.addDisplayTextBox("code", 20, 20, 100, 50, this.game.code)
    }

    processInput(input: any) {
        this.board.processInput(input)
    }

    render(ctx: CanvasRenderingContext2D) {
        this.board.render(ctx)
    }

    //socket n send

    onGameStarted(msg: ServerMessage) {
        //switch to InGame state
    }

    sendGameStart() {
        this.game.socketEventBus.emit("start-game")
    }
}

class InGame implements GameState {
    name:   string;
    board:  Board;
    world: World | null;
    game: Game;

    constructor(game: Game) {
        this.name = "in-lobby";
        this.board = new Board();
        this.world = null;
        this.game = game;
    }

    initializeBoard() {

    }

    processInput(input: any) {
        this.board.processInput(input)
        this.world?.processInput(input)
    }

    render(ctx: CanvasRenderingContext2D) {
        this.board.render(ctx)
    }

    initializeMap(msg: ServerMessage) {
        if (msg.type == "game-start") {
            let player = new Puck(msg.player.id, msg.player.position_x, msg.player.position_y)
            let opps: Puck[] = []
            msg.other_players.forEach((opp) => opps.push(new Puck(opp.id, opp.position_x, opp.position_y)))
            this.world = new World(msg.map.map, player, opps, this.game.socketEventBus)
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
    code: string;
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
            mainMenuState: new MainMenuState(this),
            requestedForLobby: new RequestedForLobby(this),
            inLobby: new InLobby(this),
            inGame: new InGame(this),
        }
        this.currentState = this.states.mainMenuState;
        this.code = "";
    }

    Play() {
    const loop = () => {
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
        const rect = this.canvas.getBoundingClientRect();
        const mousePos = {
            x: e.clientX - rect.left,
            y: e.clientY - rect.top
        };
        const button = getMouseButton(e.button);
        this.currentState.processInput({ type: "mousedown", event: e, buttonType: button, cameraX: mousePos.x, cameraY: mousePos.y });
    });

    this.canvas.addEventListener("mouseup", (e) => {
        const rect = this.canvas.getBoundingClientRect();
        const mousePos = {
            x: e.clientX - rect.left,
            y: e.clientY - rect.top
        };
        const button = getMouseButton(e.button);
        this.currentState.processInput({ type: "mouseup", event: e, buttonType: button, cameraX: mousePos.x, cameraY: mousePos.y });
    });

    this.canvas.addEventListener("mousemove", (e) => {
        const rect = this.canvas.getBoundingClientRect();
        const mousePos = {
            x: e.clientX - rect.left,
            y: e.clientY - rect.top
        };
        this.currentState.processInput({ type: "mousemove", event: e, cameraX: mousePos.x, cameraY: mousePos.y });
    });

    this.canvas.addEventListener("contextmenu", (e) => {
    e.preventDefault();
    });
}

}

function getMouseButton(button: number): MouseButton | null {
    switch (button) {
        case 0: return MouseButton.LEFT_CLICK;
        case 1: return MouseButton.MIDDLE_CLICK;
        case 2: return MouseButton.RIGHT_CLICK;
        default: return null;
    }
}