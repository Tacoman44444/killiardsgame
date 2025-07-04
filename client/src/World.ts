//this world will only be available in playstate.
//all different gamestates that can have a 'world' will be nested within this world state.

import { Arena, MapGenData } from "./Arena.js";
import { Camera } from "./camera.js";
import { SpriteComponent } from "./Components.js";
import { GameInput, MouseButton } from "./game-states.js";
import { Puck, Wall } from "./GameObjects.js";
import { Circle, Rect, ShotData, startPhysicsSimulation, Vector2 } from "./physics.js";
import { PlayerAction, ServerMessage, WallState } from "./socket-manager.js";
import { SocketEventManager } from "./socketevent-manager.js";

const radius = 16;

interface WorldState {

    name: string;
    world: World;
    Initialize() : void;
    processInput(input: any) : void;
    render(ctx: CanvasRenderingContext2D) : void;
    Enter() : void;
    Exit() : void;
}


class ActiveState implements WorldState {

    name: string;
    world: World;
    leftClickPressed: boolean;
    leftClickCoordinates: Vector2;
    constructor(world: World) {
        this.name = "active-state";
        this.world = world;
        this.leftClickPressed = false;
        this.leftClickCoordinates = new Vector2(-1, -1);
    }

    Initialize() {
        this.sub()
    }

    processInput(input: GameInput) {
        //here, we will have to process the input shot and wall placements
        if (input.type == "mousedown" && input.buttonType == MouseButton.LEFT_CLICK) {
            this.leftClickCoordinates.x = input.cameraX;
            this.leftClickCoordinates.y = input.cameraY;
            this.leftClickPressed = true;
        } else if (input.type == "mouseup" && input.buttonType == MouseButton.LEFT_CLICK) {
            if (this.leftClickPressed) {
                let directionX = this.leftClickCoordinates.x - input.cameraX;
                let directionY = this.leftClickCoordinates.y - input.cameraY;
                let action: PlayerAction = {
                    power: 1500,        // add scalable powers later
                    direction_horizontal: directionX,
                    direction_vertical: directionY,
                }
                this.sendMove(action)
            }
        } else if (input.type == "mousedown" && input.buttonType == MouseButton.RIGHT_CLICK) {
            let worldCoords = screenToWorld(input.cameraX, input.cameraY, this.world.camera, 200, 200, 64);
            if (worldCoords != null) {
                let tileCoords = worldToTile(worldCoords.x, worldCoords.y, 64);
                let wall: WallState = {
                    position_x: tileCoords.x * 64,
                    position_y: tileCoords.y * 64,
                }
                this.sendWalls(wall);
            }
        }
    }

    render(ctx: CanvasRenderingContext2D) {
        this.world.arena.render(ctx, this.world.camera);

        this.world.walls.forEach((wall) => wall.render(ctx, this.world.camera))

        this.world.opps.forEach((opp) => opp.render(ctx, this.world.camera))

        this.world.player.render(ctx, this.world.camera)
    }

    Enter() {
        console.log("entering active state")
        this.world.camera.SwitchFollow(this.world.player.circle.center);
    }

    Exit() {
        //this.unsub();
    }

    sub() {
        this.world.socketEventBus.subscribe("turn-timeout", this.onTurnTimeout.bind(this));
        this.world.socketEventBus.subscribe("wall-update", this.onWallUpdate.bind(this));
    }

    private unsub() {
        this.world.socketEventBus.unsubscribe("turn-timeout", this.onTurnTimeout);
        this.world.socketEventBus.unsubscribe("wall-update", this.onWallUpdate)
    }

    sendMove(action: PlayerAction) {
        //tell socket manager to send move
        console.log("shot vector:", action.direction_horizontal, action.direction_vertical);

        console.log("sending a move to the server");
        console.log(action)
        this.world.socketEventBus.emit("send-turn", action);
        this.world.SetState("processing-state");

        //call physics resolver
        let allCircles: Circle[] = [];
        this.world.opps.forEach((opp) => allCircles.push(opp.circle));
        allCircles.push(this.world.player.circle)
        let wallRects: Rect[] = [];
        this.world.walls.forEach((wall) => wallRects.push(wall.rect));
        console.log("starting the physics simulation")
        this.world.simInProgress = true;
        startPhysicsSimulation(this.world.player.circle, allCircles, new ShotData(new Vector2(action.direction_horizontal, action.direction_vertical), action.power), wallRects, () => {
            console.log("sending a sim-done message for my own move")
            this.world.socketEventBus.emit("simulation-done");
            if (this.world.bufferedEntityUpdate) {
                this.world.UpdateEntity(this.world.bufferedEntityUpdate)
            }
        });
    }

    sendWalls(wall: WallState) {
        this.world.socketEventBus.emit("send-wall", wall)
    }

    onTurnTimeout(msg: ServerMessage) {
        console.log("received a turn timeout message")
        this.world.SetState("processing-state")
    }

    onWallUpdate(msg: ServerMessage) {
        console.log("recieved a wall update message")
        if (this.world.currentState.name == "active-state") {   //bandage solution
            this.world.walls = [];
            if (msg.type == "wall-update") {
                msg.walls.forEach((wallState) => this.world.walls.push(new Wall(wallState.position_x, wallState.position_y)))
            }
        }
    }
}

class ProcessingState implements WorldState {

    name: string;
    world: World;
    constructor(world: World) {
        this.name = "processing-state";
        this.world = world;
    }

    Initialize() {
        console.log("subsribing")
        this.sub()
    }

    processInput(input: any) {}

    render(ctx: CanvasRenderingContext2D) {
        this.world.arena.render(ctx, this.world.camera);

        this.world.walls.forEach((wall) => wall.render(ctx, this.world.camera))

        this.world.opps.forEach((opp) => opp.render(ctx, this.world.camera))

        this.world.player.render(ctx, this.world.camera)
    }

    Enter() {
        console.log("entering processing state")
    }

    Exit() {
        //this.unsub();
    }

    sub() {
        this.world.socketEventBus.subscribe("broadcast-turn", this.onBroadcastMove.bind(this));
        this.world.socketEventBus.subscribe("entity-update", this.onEntityUpdate.bind(this));
        this.world.socketEventBus.subscribe("turn-started", this.onTurnStart.bind(this));
        this.world.socketEventBus.subscribe("wall-update", this.onWallUpdate.bind(this));
        this.world.socketEventBus.subscribe("map-update", this.onMapUpdate.bind(this));
    }

    onBroadcastMove(msg: ServerMessage) {
        //sim that physics
        console.log("recieved a move from the server");
        if (msg.type == "broadcast-turn") {
            let allCircles: Circle[] = [];
            let activeCircle: Circle | null = null;
            this.world.opps.forEach((opp) => {
                if (opp.id == msg.player.id) {
                    activeCircle = opp.circle;
                }
                allCircles.push(opp.circle);
            });
            allCircles.push(this.world.player.circle);
            let wallRects: Rect[] = [];
            this.world.walls.forEach((wall) => wallRects.push(wall.rect));
            if (activeCircle != null) {
                this.world.simInProgress = true;
                startPhysicsSimulation(activeCircle, allCircles, new ShotData(new Vector2(msg.action.direction_horizontal, msg.action.direction_vertical), msg.action.power), wallRects, () => {
                    console.log("sending a sim-done message for someone elses move")
                    this.world.simInProgress = false;
                    this.world.socketEventBus.emit("simulation-done");
                    if (this.world.bufferedEntityUpdate) {
                        this.world.UpdateEntity(this.world.bufferedEntityUpdate)
                    }
                });
            }
            
        }
    }

    onEntityUpdate(msg: ServerMessage) {

        if (this.world.simInProgress) {
            console.log("Simulation in progress. Buffering entity update.");
            this.world.bufferedEntityUpdate = msg;
            return;
        }

        if (msg.type == "entity-update") {
            console.log("recieved an entity update message")
            let allPlayers: Puck[] = [...this.world.opps, this.world.player]
            allPlayers.sort((a, b) => a.id.localeCompare(b.id, undefined, { numeric: true }));
            msg.all_players.sort((a, b) => a.id.localeCompare(b.id, undefined, {numeric: true}));

            for (let i = 0; i < allPlayers.length; i++) {
                allPlayers[i].circle.update(msg.all_players[i].position_x, msg.all_players[i].position_y, msg.all_players[i].velocity_x, msg.all_players[i].velocity_y)
            }

            this.world.walls = [];
            msg.walls.forEach((wallState) => this.world.walls.push(new Wall(wallState.position_x, wallState.position_y)))
        }
    }

    onTurnStart(msg: ServerMessage) {
        console.log("recieved a turn start message")
        this.world.SetState("active-state")
    }

    onWallUpdate(msg: ServerMessage) {
        if (this.world.currentState.name == "processing-state") {
            if (msg.type == "wall-update") {
            this.world.walls = [];
            msg.walls.forEach((wallState) => this.world.walls.push(new Wall(wallState.position_x, wallState.position_y)))
            }
        }

    }

    onMapUpdate(msg: ServerMessage) {
        if (msg.type == "map-update") {

        }
    }
}

export class World {

    player: Puck;
    opps: Puck[] = [];
    arena: Arena;
    walls: Wall[] = [];
    camera: Camera;
    states: {
        activeState: WorldState;
        processingState: WorldState;
    }
    currentState: WorldState;
    socketEventBus: SocketEventManager
    simInProgress: boolean = false;
    bufferedEntityUpdate: ServerMessage | null = null;

    constructor(mapData: MapGenData, player: Puck, opps: Puck[], socketEventBus: SocketEventManager) {
        this.player = player;
        this.opps = opps
        this.arena = new Arena(new SpriteComponent("floorTile"), new SpriteComponent("abyssTile"), mapData);
        this.camera = new Camera(player.circle.center, 1600, 1200);
        this.states = {
            activeState: new ActiveState(this),
            processingState: new ProcessingState(this),
        }
        this.currentState = this.states.processingState;
        this.socketEventBus = socketEventBus;
        this.currentState.Enter();
        this.states.activeState.Initialize();
        this.states.processingState.Initialize();
    }

    UpdateEntity(msg: ServerMessage) {
        if (msg.type == "entity-update") {
            console.log("this is a delayed entity update rahhh")
            let allPlayers: Puck[] = [...this.opps, this.player];
            console.log("client length: ", allPlayers.length)
            console.log("server length: ", msg.all_players.length)
            allPlayers.sort((a, b) => a.id.localeCompare(b.id, undefined, { numeric: true }));
            msg.all_players.sort((a, b) => a.id.localeCompare(b.id, undefined, {numeric: true}));

            for (let i = 0; i < allPlayers.length; i++) {
                allPlayers[i].circle.update(msg.all_players[i].position_x, msg.all_players[i].position_y, msg.all_players[i].velocity_x, msg.all_players[i].velocity_y)
            }
            this.walls = [];
            msg.walls.forEach((wallState) => this.walls.push(new Wall(wallState.position_x, wallState.position_y)))
        }
    }

    SetState(state: string) {
        switch (state) {
            case "active-state":
                this.currentState.Exit();
                this.currentState = this.states.activeState;
                this.currentState.Enter();
                break;
            case "processing-state":
                this.currentState.Exit();
                this.currentState = this.states.processingState;
                this.currentState.Enter();
                break;
            default:
                console.log("invalid state");
                break;
        }
    }

    processInput(input: any) {
        this.currentState.processInput(input);
    }

    render(ctx: CanvasRenderingContext2D) {
        this.currentState.render(ctx);
    }

}

function screenToWorld(mouseX: number, mouseY: number, camera: Camera, mapWidth: number, mapHeight: number, tileSize: number): Vector2 | null {
    const worldX = mouseX + (camera.follow.x - camera.width / 2);
    const worldY = mouseY + (camera.follow.y - camera.height / 2);

    const maxWorldX = mapWidth * tileSize;
    const maxWorldY = mapHeight * tileSize;

    if (worldX < 0 || worldY < 0 || worldX >= maxWorldX || worldY >= maxWorldY) {
        return null; // Out of bounds
    }

    return new Vector2(worldX, worldY);
}

function worldToTile(worldX: number, worldY: number, tileSize: number): Vector2 {
    return new Vector2(Math.floor(worldX / tileSize), Math.floor(worldY / tileSize));
}