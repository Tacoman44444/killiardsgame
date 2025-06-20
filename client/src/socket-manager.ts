//the primary role of the socket manager is to communicate with the server

import { MapGenData } from "./Arena";

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
        id: string;
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
        state: WallState;
    }
    | {
        type: "send-turn";
        action: PlayerAction;
    };

type ServerMessage = 
    | {
        type: "room-created";
        code: string;
    }
    | {
        type: "game-start";
        map: MapState;
        playerPos: PlayerState;
        otherPlayers: PlayerState[];

    }
    | {
        type: "turn-start";
    }
    | {
        type: "turn-timeout";
    }
    | {
        type: "player-update";
        playerId: string;
        newState: PlayerState;
    }
    | {
        type: "map-update";
        map: MapState;
    }
    | {
        type: "game-finished";
        winnerName: string;
    }
    | {
        type: "wall-update";
        walls: WallState[];
    };


class SocketManager {
    private socket: WebSocket;
    
    constructor(url: string) {
        this.socket = new WebSocket(url);
    }

    private send(msg: ClientMessage) {
        this.socket.send(JSON.stringify(msg));
    }

    sendCreateRoomRequest(playerID: string) {
        const msg: ClientMessage = {
            type: "create-room",
            id: playerID,
        };
        this.send(msg);
    }

    sendStartGameRequest() {
        const msg: ClientMessage = {
            type: "start-game",
        };
        this.send(msg);
    }

    sendJoinRoomRequest(roomCode: string) {
        const msg: ClientMessage = {
            type: "join-room",
            code: roomCode,
        };
        this.send(msg);
    }

    sendWall(wallState: WallState) {
        const msg: ClientMessage = {
            type: "send-wall",
            state: wallState,
        };
        this.send(msg);
    }

    sendTurn(PlayerAction: PlayerAction) {
        const msg: ClientMessage = {
            type: "send-turn",
            action: PlayerAction,
        }
    }
}