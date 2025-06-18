//this world will only be available in playstate.
//all different gamestates that can have a 'world' will be nested within this world state.

import { Arena } from "./Arena";
import { SpriteComponent } from "./Components";
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

    arena: Arena;
    states: {
        activeState: WorldState;
        processingState: WorldState;
        inactiveState: WorldState;
    }
    currentState: WorldState;

    constructor() {
        this.arena = new Arena(new SpriteComponent("floorTile"), new SpriteComponent("abyssTile"));
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