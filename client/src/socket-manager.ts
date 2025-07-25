//the primary role of the socket manager is to communicate with the server

import { MapGenData } from "./Arena.js";
import { SocketEventManager } from "./socketevent-manager.js";

export type PlayerIdentity = {
    id:         string;
    position_x: number;
    position_y: number;
    velocity_x: number;
    velocity_y: number;
    username:   string;
}

export type PlayerAction = {
    power: number;
    direction_horizontal: number;
    direction_vertical: number;
}

export type MapState = {
    arena: number[][],
    width: number,
    height: number,
}

export type WallState = {
    position_x: number;
    position_y: number;
}

export type JoinRoomData = {
    username: string,
    code: string
}

type ClientMessage = 
    | {
        type: "create-room";
        username: string,
    }
    | {
        type: "start-game";
    }
    | {
        type: "join-room";
        data: JoinRoomData;
    }
    | {
        type: "leave-room";
    }
    | {
        type: "send-wall";
        wall_state: WallState;
    }
    | {
        type: "send-turn";
        player_action: PlayerAction;
    } 
    | {
        type: "simulation-done"
    }
    | {
        type: "return-to-mainmenu"
    }
    | {
        type: "return-to-lobby"
    };

export type ServerMessage = 
    | {
        type: "room-created";
        code: string;
    }
    | {
        type: "room-joined";
        code: string;
    }
    | {
        type: "invalid-code";
    }
    | {
        type: "game-start";
        current_map: MapState;
        next_map: MapState
        player: PlayerIdentity;
        other_players: PlayerIdentity[];
    }
    | {
        type: "turn-start";
    }
    | {
        type: "turn-timeout";
    }
    | {
        type: "broadcast-turn";
        player: PlayerIdentity;
        action: PlayerAction;
    }
    | {
        type: "entity-update";
        player: PlayerIdentity;
        all_players: PlayerIdentity[];
        walls: WallState[]
    }
    | {
        type: "e"
        eliminated_players: PlayerIdentity[]
    }
    | {
        type: "wall-update";
        walls: WallState[];
    }
    | {
        type: "map-update"; 
        current_map: MapState;
        next_map: MapState;
    }
    | {
        type: "game-finished";
        result: string;
        winner_name: string;
    }
    | {
        type: "lobby-closed";
    }
    | {
        type: "return-to-lobby";
    };


export class SocketManager {
    private socket: WebSocket;
    private eventManager: SocketEventManager;
    
    constructor(url: string, eventManager: SocketEventManager) {
        this.socket = new WebSocket(url);
        this.eventManager = eventManager;
        this.socket.onerror = (e) => console.error("WebSocket error", e);
        this.socket.onmessage = (e) => this.handleMessage(e);
        this.subscribeToEventBus();
    }

    async connect(): Promise<void> {
        return new Promise((resolve, reject) => {
            this.socket.onopen = () => {
                console.log("connected");
                resolve();
            };
            this.socket.onerror = (e) => {
                reject(e);
            };
        });
    }

    private send(msg: ClientMessage) {
        this.socket.send(JSON.stringify(msg));
    }

    sendCreateRoomRequest(name: string) {
        console.log("sending create room request")
        const msg: ClientMessage = {
            type: "create-room",
            username: name,
        };
        this.send(msg);
    }

    sendStartGameRequest() {
        const msg: ClientMessage = {
            type: "start-game",
        };
        this.send(msg);
    }

    sendJoinRoomRequest(data: JoinRoomData) {
        console.log("sending join room request")
        const msg: ClientMessage = {
            type: "join-room",
            data: data,
        };
        this.send(msg);
    }

    sendLeaveRoomRequest() {
        const msg: ClientMessage = {
            type: "leave-room",
        };
        this.send(msg)
    }

    sendWall(wallState: WallState) {
        const msg: ClientMessage = {
            type: "send-wall",
            wall_state: wallState,
        };
        this.send(msg);
    }

    sendTurn(PlayerAction: PlayerAction) {
        const msg: ClientMessage = {
            type: "send-turn",
            player_action: PlayerAction,
        }
        this.send(msg);
    }

    sendSimulationDone() {
        const msg: ClientMessage = {
            type: "simulation-done"
        };
        this.send(msg);
    }

    sendReturnToMainMenu() {
        const msg: ClientMessage = {
            type: "return-to-mainmenu"
        };
        this.send(msg)
    }

    sendReturnToLobby() {
        const msg: ClientMessage = {
            type: "return-to-lobby"
        };
        this.send(msg)
    }

    private subscribeToEventBus() {
        this.eventManager.subscribe("create-room", this.sendCreateRoomRequest.bind(this));
        this.eventManager.subscribe("start-game", this.sendStartGameRequest.bind(this));
        this.eventManager.subscribe("join-room", this.sendJoinRoomRequest.bind(this));
        this.eventManager.subscribe("leave-room", this.sendLeaveRoomRequest.bind(this));
        this.eventManager.subscribe("send-wall", this.sendWall.bind(this));
        this.eventManager.subscribe("send-turn", this.sendTurn.bind(this));
        this.eventManager.subscribe("simulation-done", this.sendSimulationDone.bind(this));
        this.eventManager.subscribe("return-to-mainmenu", this.sendReturnToMainMenu.bind(this));
        this.eventManager.subscribe("return-to-lobby", this.sendReturnToLobby.bind(this));
    }

    private handleMessage(e: MessageEvent) {
        const raw = e.data as string;
        let msg: ServerMessage;
        try {
            msg = JSON.parse(raw) as ServerMessage;
        } catch (err) {
            console.error("Bad JSON from the server: ", err);
            return;
        }
        if (msg.type == "broadcast-turn") {
            console.log("ACTUALLY recieved a move from the server")
        }
        this.eventManager.emit(msg.type, msg);
    }
}