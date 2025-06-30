//this world will only be available in playstate.
//all different gamestates that can have a 'world' will be nested within this world state.

import { Arena, MapGenData } from "./Arena";
import { Camera } from "./camera";
import { SpriteComponent } from "./Components";
import { Puck, Wall } from "./GameObjects";
import { Circle, physicsResolver, ShotData, Vector2 } from "./physics";
import { PlayerAction, ServerMessage } from "./socket-manager";
import { SocketEventManager } from "./socketevent-manager";

const radius = 16;

interface WorldState {

    name: string;
    world: World;
    processInput(input: any) : void;
    update() : void;
    render(ctx: CanvasRenderingContext2D) : void;
    Enter() : void;
}

class ActiveState implements WorldState {

    name: string;
    world: World;
    constructor(world: World) {
        this.name = "active_state";
        this.world = world;
    }

    processInput(input: any) {
        //here, we will have to process the input shot and wall placements
        
    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D) {
        this.world.arena.render(ctx, this.world.camera);

        this.world.walls.forEach((wall) => wall.render(ctx, this.world.camera))

        this.world.opps.forEach((opp) => opp.render(ctx, this.world.camera))

        this.world.player.render(ctx, this.world.camera)
    }

    Enter() {
        this.world.camera.SwitchFollow(this.world.player.position);
    }

    sendMove(action: PlayerAction) {
        //tell socket manager to send move
        this.world.socketEventBus.emit("send-turn", action)
        this.world.SetState("processing-state")

        //call physics resolver
    }

    onTurnTimeout(msg: ServerMessage) {
        this.world.SetState("processing-state")
    }

    onWallUpdate(msg: ServerMessage) {
        this.world.walls = [];
        if (msg.type == "wall-update") {
            msg.walls.forEach((wallState) => this.world.walls.push(new Wall(wallState.position_x, wallState.position_y)))
        }
    }

    
}

class ProcessingState implements WorldState {

    name: string;
    world: World;
    constructor(world: World) {
        this.name = "processing_state";
        this.world = world;
    }

    processInput(input: any) {
        //here, there is nothing to process for now

    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D) {

    }

    Enter() {

    }

    onEntityUpdate(msg: ServerMessage) {
        //check with your own physics simulation.
    }

    onTurnStart(msg: ServerMessage) {
        this.world.SetState("active-state")
    }

    onWallUpdate(msg: ServerMessage) {
        this.world.walls = [];
        if (msg.type == "wall-update") {
            msg.walls.forEach((wallState) => this.world.walls.push(new Wall(wallState.position_x, wallState.position_y)))
        }
    }

    onMapUpdate(msg: ServerMessage) {

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
    constructor(mapData: MapGenData, player: Puck, opps: Puck[], socketEventBus: SocketEventManager) {
        this.player = player;
        this.opps = opps
        this.arena = new Arena(new SpriteComponent("floorTile"), new SpriteComponent("abyssTile"), mapData);
        this.camera = new Camera(player.position, 1600, 1200);
        this.states = {
            activeState: new ActiveState(this),
            processingState: new ProcessingState(this),
        }
        this.currentState = this.states.processingState;
        this.socketEventBus = socketEventBus;
    }

    SetState(state: string) {
        switch (state) {
            case "active-state":
                this.currentState = this.states.activeState;
                this.currentState.Enter();
                break;
            case "processing-state":
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

    update() {
        this.currentState.update();
    }

    render(ctx: CanvasRenderingContext2D) {
        this.currentState.render(ctx);
    }

}