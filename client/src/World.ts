//this world will only be available in playstate.
//all different gamestates that can have a 'world' will be nested within this world state.

import { Arena, MapGenData } from "./Arena";
import { Camera } from "./camera";
import { SpriteComponent } from "./Components";
import { Puck, Wall } from "./GameObjects";

interface WorldState {

    name: string;
    world: World;
    processInput() : void;
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

    processInput() {

    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D) {

    }

    Enter() {

    }
}

class ProcessingState implements WorldState {

    name: string;
    world: World;
    constructor(world: World) {
        this.name = "processing_state";
        this.world = world;
    }

    processInput() {

    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D) {

    }

    Enter() {

    }
}

class InactiveState implements WorldState {

    name: string;
    world: World;
    constructor(world: World) {
        this.name = "inactive-state";
        this.world = world;
    }

    processInput() {

    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D) {

    }

    Enter() {

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
        inactiveState: WorldState;
    }
    currentState: WorldState;

    constructor(mapData: MapGenData, player: Puck, opps: Puck[]) {
        this.player = player;
        this.opps = opps
        this.arena = new Arena(new SpriteComponent("floorTile"), new SpriteComponent("abyssTile"), mapData);
        this.camera = new Camera(player.position, 1600, 1200);
        this.states = {
            activeState: new ActiveState(this),
            processingState: new ProcessingState(this),
            inactiveState: new InactiveState(this),
        }
        this.currentState = this.states.inactiveState;
    }

    processInput() {
        this.currentState.processInput();
    }

    update() {
        this.currentState.update();
    }

    render(ctx: CanvasRenderingContext2D) {
        this.currentState.render(ctx);
    }

}