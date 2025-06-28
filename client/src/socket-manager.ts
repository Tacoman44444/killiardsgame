//the primary role of the socket manager is to communicate with the server

import { MapGenData } from "./Arena";
import { SocketEventManager } from "./socketevent-manager";

type PlayerState = {
    position_x: number;
    position_y: number;
    velocity_x: number;
    velocity_y: number;
}

type PlayerAction = {
    power: number;
    direction_horizontal: number;
    direction_verical: number;
}

type MapState = {
    map: MapGenData;
}

type WallState = {
    position_x: number;
    position_y: number;
}

type ClientMessage = 
    | {
        type: "create-room";
        name: string
    }
    | {
        type: "start-game";
    }
    | {
        type: "join-room";
        name: string
        code: number;
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

type ServerMessage = 
    | {
        type: "room-created";
        code: number;
    }
    | {
        type: "room-joined"
    }
    | {
        type: "game-start";
        map: MapState;
        player: PlayerState;
        other_players: PlayerState[];
    }
    | {
        type: "turn-start";
    }
    | {
        type: "turn-timeout";
    }
    | {
        type: "entity-update";
        player: PlayerState;
        other_players: PlayerState[];
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

    sendCreateRoomRequest(player_name: string) {
        const msg: ClientMessage = {
            type: "create-room",
            name: player_name
        };
        this.send(msg);
    }

    sendStartGameRequest() {
        const msg: ClientMessage = {
            type: "start-game",
        };
        this.send(msg);
    }

    sendJoinRoomRequest(player_name: string, roomCode: number) {
        const msg: ClientMessage = {
            type: "join-room",
            name: player_name,
            code: roomCode,
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