//the primary role of the socket manager is to communicate with the server

import { MapGenData } from "./Arena";
import { SocketEventManager } from "./socketevent-manager";

export type PlayerIdentity = {
    id:         string;
    position_x: number;
    position_y: number;
    velocity_x: number;
    velocity_y: number;
}

export type PlayerAction = {
    power: number;
    direction_horizontal: number;
    direction_verical: number;
}

export type MapState = {
    map: MapGenData;
}

export type WallState = {
    position_x: number;
    position_y: number;
}

export type JoinRoomData = {
    id: string,
    code: string
}

type ClientMessage = 
    | {
        type: "create-room";
    }
    | {
        type: "start-game";
    }
    | {
        type: "join-room";
        code: string;
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
        map: MapState;
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
        other_players: PlayerIdentity[];
        walls: WallState[]
    }
    | {
        type: "wall-update";
        walls: WallState[];
    }
    | {
        type: "map-update"; 
        map: MapState;
    }
    | {
        type: "game-finished";
        winner_name: string;
    };


class SocketManager {
    private socket: WebSocket;
    private eventManager: SocketEventManager;
    
    constructor(url: string, eventManager: SocketEventManager) {
        this.socket = new WebSocket(url);
        this.eventManager = eventManager;
        this.socket.onmessage = (e) => this.eventManager;
    }

    private send(msg: ClientMessage) {
        this.socket.send(JSON.stringify(msg));
    }

    sendCreateRoomRequest() {
        const msg: ClientMessage = {
            type: "create-room",
        };
        this.send(msg);
    }

    sendStartGameRequest() {
        const msg: ClientMessage = {
            type: "start-game",
        };
        this.send(msg);
    }

    sendJoinRoomRequest(code: string) {
        const msg: ClientMessage = {
            type: "join-room",
            code: code,
        };
        this.send(msg);
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
        this.eventManager.emit(msg.type, msg);
    }
}