import { SpriteComponent } from "./Components.js"
import { Puck } from "./GameObjects.js"
import { ServerMessage } from "./socket-manager.js"
import { SocketEventManager } from "./socketevent-manager.js"
import { SoundManager } from "./sound-manager.js"
import { Board, BoardEventManager, Img } from "./ui.js"
import { POWER_LEVEL, World } from "./World.js"


interface GameState {
    name:       string
    board:      Board
    boardEventManager: BoardEventManager
    game: Game
    Enter() : void;
    initializeBoard() : void;
    initializeWorld(world: World) : void;
    processInput(input: any) : void;
    render(ctx: CanvasRenderingContext2D) : void;
    Exit() : void;
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
    boardEventManager: BoardEventManager
    game: Game;

    constructor(game: Game) {
        this.name = "main-menu";
        this.board = new Board();
        this.boardEventManager = new BoardEventManager(this.board)
        this.game = game;
        this.initializeBoard();
    }

    Enter() {

    }

    initializeBoard() {
        this.board.addInputTextBox("username", 650, 400, 300, 50, 13);
        this.board.addInputTextBox("code", 725, 500, 150, 40, 6);
        this.board.addTextButton("create", 290, 700, 224, 48, "CREATE ROOM", () => this.sendCreateRoomRequest(), true, 24);
        this.board.addTextButton("join", 1100, 700, 204, 48, "JOIN ROOM", () => this.sendJoinRequest(), true, 24)
        this.board.addDisplayTextBox("test", 400, 100, 100, 100, "KILLIARDS", 128, true)

    }

    initializeWorld(world: World): void {}

    processInput(input: any) {
        this.board.processInput(input)
    }

    render(ctx: CanvasRenderingContext2D) {
        console.log("in the main menu render loop");
        this.board.render(ctx)
    }

    getName(): string {
        return this.board.inputBoxes.find(box => box.name === "username")?.text || "";
    }

    getCode(): string {
        return this.board.inputBoxes.find(box => box.name === "code")?.text || "";
    }

    Exit() {

    }

    //handle socket event and send messages now ... 

    sendJoinRequest() {
        let code = this.getCode()
        let name = this.getName()

        if (code == "" || name == "") {
            console.log("invalid details")
        } else {
            console.log("sending this code to the server: ", code)
            this.game.socketEventBus.emit("join-room", { code: code, name: name });
            this.game.SetState(this.game.states.requestedForLobby)
        }

    }

    sendCreateRoomRequest() {
        let name = this.getName()
        if (name === "") {
            console.log("invalid details")
        }
        this.game.socketEventBus.emit("create-room", name)
        this.game.SetState(this.game.states.requestedForLobby)
    }
}

class RequestedForLobby implements GameState {
    name:   string;
    board:  Board;
    boardEventManager: BoardEventManager
    game: Game;
    code: string;   //rage coding, probably shouldnt do it this way

    constructor(game: Game) {
        this.name = "requested-for-lobby";
        this.board = new Board();
        this.boardEventManager = new BoardEventManager(this.board)
        this.game = game;
        this.code = "";
        this.sub();
        this.initializeBoard();
    }

    Enter() {

    }

    initializeBoard() {
        
    }

    initializeWorld(world: World): void {}

    processInput(input: any) {
        this.board.processInput(input)
    }

    render(ctx: CanvasRenderingContext2D) {
        this.board.render(ctx)
        console.log("in the requested for lobby render loop");
    }

    private sub() {
        this.game.socketEventBus.subscribe("invalid-code", this.onInvalidCode.bind(this));
        this.game.socketEventBus.subscribe("room-created", this.onRoomCreated.bind(this))
        this.game.socketEventBus.subscribe("room-joined", this.onRoomJoined.bind(this));
    }

    Exit() {

    }

    //socket msges and send

    onInvalidCode(msg: ServerMessage) {
        this.game.SetState(this.game.states.mainMenuState)
    }

    onRoomCreated(msg: ServerMessage) {
        if (msg.type == "room-created") {
            this.game.SetState(this.game.states.inLobby)
            this.game.currentState.boardEventManager.onEvent("roomcodeupdate", msg.code); 
        }
        
    }

    onRoomJoined(msg: ServerMessage) {
        if (msg.type == "room-joined") {
            this.game.SetState(this.game.states.inLobby)
            this.game.currentState.boardEventManager.onEvent("roomcodeupdate", msg.code);
        }

    }
}

class InLobby implements GameState {
    name:   string;
    board:  Board;
    boardEventManager: BoardEventManager
    game: Game;

    constructor(game: Game) {
        this.name = "in-lobby";
        this.board = new Board();
        this.boardEventManager = new BoardEventManager(this.board)
        this.game = game;
        this.initializeBoard();
        this.sub();
    }

    Enter() {

    }

    initializeBoard() {
        console.log("the code is: ", this.game.code)
        this.board.addDisplayTextBox("code", 150, 20, 100, 50, this.game.code, 48, true)
        this.board.addTextButton("start", 700, 700, 250, 60, "START", () => this.sendGameStart(), true, 48)
        this.board.addTextButton("leave", 700, 800, 250, 60, "LEAVE", () => this.sendLeaveRoom(), true, 48)

        this.boardEventManager.addEvent("roomcodeupdate", (code: string) => {
            this.boardEventManager.board.getDisplayTextBox("code")?.UpdateText(code)
        })
    }

    initializeWorld(world: World): void {}

    processInput(input: any) {
        this.board.processInput(input)
    }

    render(ctx: CanvasRenderingContext2D) {
        this.board.render(ctx)
        console.log("in the in lobby render loop");
    }

    private sub() {
        this.game.socketEventBus.subscribe("game-start", this.onGameStarted.bind(this));
    }

    Exit() {

    }

    //socket n send

    onGameStarted(msg: ServerMessage) {
        if (msg.type == "game-start") {
            console.log("GAME START RECEIVED, msg =", JSON.stringify(msg, null, 2));
            this.game.SetState(this.game.states.inGame)
            let player = new Puck(msg.player.id, msg.player.username, msg.player.position_x, msg.player.position_y)
            let opps: Puck[] = []
            msg.other_players.forEach((opp) => opps.push(new Puck(opp.id, opp.username, opp.position_x, opp.position_y)))
            let world = new World(msg.current_map, msg.next_map, player, opps, this.game.socketEventBus, this.game.states.inGame.boardEventManager, this.game.soundManager)
            this.game.currentState.initializeWorld(world);
        }
    }

    sendGameStart() {
        this.game.socketEventBus.emit("start-game")
        console.log("sending game start")
    }

    sendLeaveRoom() {
        this.game.socketEventBus.emit("leave-room")
        this.game.SetState(this.game.states.mainMenuState)
    }
}

class InGame implements GameState {
    name:   string;
    board:  Board;
    boardEventManager: BoardEventManager;
    world: World | null;
    game: Game;

    constructor(game: Game) {
        this.name = "in-lobby";
        this.board = new Board();
        this.boardEventManager = new BoardEventManager(this.board);
        this.world = null;
        this.game = game;
        this.initializeBoard()
        this.sub()
    }

    Enter() {
        let winBox = this.board.getDisplayTextBox("win")
        if (winBox) { winBox.visible = false }

        let drawBox = this.board.getDisplayTextBox("draw")
        if (drawBox) { drawBox.visible = false }

        let gameoverBox = this.board.getDisplayTextBox("gameover")
        if (gameoverBox) { gameoverBox.visible = false }

        let lobbyBtn = this.board.getTextButton("return-to-lobby")
        if (lobbyBtn) { lobbyBtn.visible = false }

        let menuBtn = this.board.getTextButton("return-to-mainmenu")
        if (menuBtn) { menuBtn.visible = false }
    }

    private sub() {
        this.game.socketEventBus.subscribe("lobby-closed", this.onLobbyClosed.bind(this))
        this.game.socketEventBus.subscribe("return-to-lobby", this.onReturnToLobby.bind(this))
    }

    initializeBoard() {
        this.board.addDisplayTextBox("spectating", 20, 20, 200, 100, "", 48, false);
        this.board.addDisplayTextBox("turnactive", 400, 100, 1000, 200, "YOUR TURN -- SHOOT", 48, false);
        this.board.addImage("1", 40, 1060, 50, 100, new SpriteComponent("power_level_1"), false);
        this.board.addImage("2", 40, 950, 50, 100, new SpriteComponent("power_level_2"), false);
        this.board.addImage("3", 40, 840, 50, 100, new SpriteComponent("power_level_3"), false);
        this.board.addImage("4", 40, 730, 50, 100, new SpriteComponent("power_level_4"), false);
        this.board.addImage("5", 40, 620, 50, 100, new SpriteComponent("power_level_5"), false);
        this.board.addDisplayTextBox("gameover", 800, 250, 1200, 400, "GAME OVER", 98, false)
        this.board.addDisplayTextBox("win", 800, 500, 200, 40, "", 64, false)
        this.board.addDisplayTextBox("draw", 800, 500, 200, 40, "DRAW", 64, false)
        this.board.addTextButton("return-to-lobby", 800, 700, 250, 60, "LOBBY", () => this.sendReturnToLobby(), false, 48)
        this.board.addTextButton("return-to-mainmenu", 800, 800, 250, 60, "HUB", () => this.sendReturnToMainMenu(), false, 48)


        this.boardEventManager.addEvent("winner", (name: string) => {
            let gameoverBox = this.boardEventManager.board.getDisplayTextBox("gameover")
            if (gameoverBox) {
                gameoverBox.visible = true;
            }
            let winBox = this.boardEventManager.board.getDisplayTextBox("win")
            if (winBox) {
                winBox.visible = true;
                winBox.UpdateText("WINNER: " + name)
            }
        })
        this.boardEventManager.addEvent("draw", () => {
            let gameoverBox = this.boardEventManager.board.getDisplayTextBox("gameover")
            if (gameoverBox) {
                gameoverBox.visible = true;
            }

            let drawBox = this.boardEventManager.board.getDisplayTextBox("draw")
            if (drawBox) {
                drawBox.visible = true;
            }
        })
        this.boardEventManager.addEvent("spectating", (username: string) => {
            let box = this.board.getDisplayTextBox("spectating");
            if (box) {
                box.visible = true;
                box.text = "SPECTATING " + username;
            } else {
                console.log("WHERE IS THE BOX")
            }
        })
        this.boardEventManager.addEvent("onplayer", () => {
            let box = this.board.getDisplayTextBox("spectating");
            if (box) {
                box.visible = false;
            } else {
                console.log("WHERE IS THE BOX")
            }
        })
        this.boardEventManager.addEvent("turnactive", () => {
            let box = this.board.getDisplayTextBox("turnactive");
            if (box) {
                box.visible = true;
            } else {
                console.log("WHERE IS THE BOX")
            }
            let boxspec = this.board.getDisplayTextBox("spectating");
            if (boxspec) {
                boxspec.visible = false;
            } else {
                console.log("WHERE IS THE BOX")
            }
        })
        this.boardEventManager.addEvent("turncompleted", () => {
            let box = this.board.getDisplayTextBox("turnactive");
            if (box) {
                box.visible = false;
            } else {
                console.log("WHERE IS THE BOX")
            }
        })
        this.boardEventManager.addEvent("showpowerlevel", (powerLevel: POWER_LEVEL) => {
            let images: Img[] = []
            for (const img of this.board.images) {
                if (!isNaN(Number(img.name)) && img.name.trim() !== "") {
                    const number = Number(img.name);
                    if (number <= powerLevel) {
                        img.visible = true;
                    } else {
                        img.visible = false;
                    }
                }
            }
            
        })
        this.boardEventManager.addEvent("hidepowerlevel", () => {
            let images: Img[] = [] 
            for (const img of this.board.images) {
                if (!isNaN(Number(img.name)) && img.name.trim() !== "") {
                    img.visible = false
                }
            }
        })
        this.boardEventManager.addEvent("showbuttons", () => {
            let lobbyBtn = this.board.getTextButton("return-to-lobby");
            if (lobbyBtn) {
                lobbyBtn.visible = true;
            }

            let menuBtn = this.board.getTextButton("return-to-mainmenu");
            if (menuBtn) {
                menuBtn.visible = true;
            }
        })
    }

    initializeWorld(world: World) {
        this.world = world;
        console.log("initialized world");
    }

    processInput(input: any) {
        this.board.processInput(input)
        this.world?.processInput(input)
    }

    render(ctx: CanvasRenderingContext2D) {
        if (this.world != null) {
            this.world.render(ctx);
        }

        this.board.render(ctx); 
    }

    Exit() {
        
    }

    //handle them socket events and send them messages

    sendReturnToLobby() {
        this.game.socketEventBus.emit("return-to-lobby")
    }

    sendReturnToMainMenu() {
        this.game.socketEventBus.emit("return-to-mainmenu")
        this.game.SetState(this.game.states.mainMenuState)
    }

    onLobbyClosed(msg: ServerMessage) {
        if (msg.type == "lobby-closed") {
            this.game.SetState(this.game.states.mainMenuState)
        }
    }

    onReturnToLobby(msg: ServerMessage) {
        if (msg.type == "return-to-lobby") {
            this.game.SetState(this.game.states.inLobby)
        }
    }    

}

export class Game {
    canvas: HTMLCanvasElement;
    ctx: CanvasRenderingContext2D;
    currentState: GameState;
    socketEventBus: SocketEventManager;
    soundManager: SoundManager;
    code: string;
    scale: number;
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
        this.soundManager = new SoundManager()
        this.states = {
            mainMenuState: new MainMenuState(this),
            requestedForLobby: new RequestedForLobby(this),
            inLobby: new InLobby(this),
            inGame: new InGame(this),
        }
        this.currentState = this.states.mainMenuState;
        this.code = "";
        this.scale = 0;
        this.resizeWindow();
    }

    SetState(state: GameState) {
        this.currentState.Exit();
        this.currentState = state;
        this.currentState.Enter();
    }

    Play() {
    const loop = () => {
        this.ctx.setTransform(1, 0, 0, 1, 0, 0);
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

        this.ctx.setTransform(this.scale, 0, 0, this.scale, 0, 0); // apply scale

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
            x: (e.clientX - rect.left) / this.scale,
            y: (e.clientY - rect.top) / this.scale
        };
        const button = getMouseButton(e.button);
        this.currentState.processInput({ type: "mousedown", event: e, buttonType: button, cameraX: mousePos.x, cameraY: mousePos.y });
    });

    this.canvas.addEventListener("mouseup", (e) => {
        const rect = this.canvas.getBoundingClientRect();
        const mousePos = {
            x: (e.clientX - rect.left) / this.scale,
            y: (e.clientY - rect.top) / this.scale
        };
        const button = getMouseButton(e.button);
        this.currentState.processInput({ type: "mouseup", event: e, buttonType: button, cameraX: mousePos.x, cameraY: mousePos.y });
    });

    this.canvas.addEventListener("mousemove", (e) => {
        const rect = this.canvas.getBoundingClientRect();
        const mousePos = {
            x: (e.clientX - rect.left) / this.scale,
            y: (e.clientY - rect.top) / this.scale
        };
        this.currentState.processInput({ type: "mousemove", event: e, cameraX: mousePos.x, cameraY: mousePos.y });
    });

    this.canvas.addEventListener("contextmenu", (e) => {
        e.preventDefault();
    });

    window.addEventListener("resize", () => {
        this.resizeWindow();
    });
    }

    resizeWindow() {
        const windowWidth = window.innerWidth;
        const windowHeight = window.innerHeight;

        // Target aspect ratio: 4:3
        let newWidth = windowWidth;
        let newHeight = (newWidth * 3) / 4;

        if (newHeight > windowHeight) {
            newHeight = windowHeight;
            newWidth = (newHeight * 4) / 3;
        }

        // Set actual canvas drawing size (resolution)
        this.canvas.width = newWidth;
        this.canvas.height = newHeight;

        // Set CSS size (how big it looks)
        this.canvas.style.width = `${newWidth}px`;
        this.canvas.style.height = `${newHeight}px`;

        const scaleX = this.canvas.width / 1600;
        const scaleY = this.canvas.height / 1200;
        this.scale = Math.min(scaleX, scaleY); 
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