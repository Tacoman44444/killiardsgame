//all gameobject implementations

import { Camera } from "./camera";
import { PositionComponent, SpriteComponent, VelocityComponent } from "./Components";
import { Circle, Vector2 } from "./physics";

interface GameObject {

    name: string;

    processInput() : void;
    update() : void;
    render(ctx : CanvasRenderingContext2D, camera: Camera) : void;
}

export class Puck implements GameObject {
    
    name: string;
    puckSprite: SpriteComponent;
    position: PositionComponent;
    velocity: VelocityComponent;

    constructor(posX: number, posY: number) {
        this.name = "puck";
        this.puckSprite = new SpriteComponent("puck");
        this.position = new PositionComponent(posX, posY);
        this.velocity = new VelocityComponent(0, 0);
    }

    processInput() {

    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D, camera: Camera) {

    }
}

class Brick implements GameObject {

    name: string;
    woodSprite: SpriteComponent;
    position: PositionComponent;

    constructor(posX: number, posY: number) {
        this.name = "brick"
        this.woodSprite = new SpriteComponent("brickTile");
        this.position = new PositionComponent(posX, posY);
    }

    processInput() {

    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D, camera: Camera) {

    }
}

export class Wall implements GameObject {

    name: string;
    wallSprite: SpriteComponent;
    position: PositionComponent;

    constructor(posX: number, posY: number) {
        this.name = "wall";
        this.wallSprite = new SpriteComponent("wallTile");
        this.position = new PositionComponent(posX, posY);
    }

    processInput() {

    }

    update() {

    }

    render(ctx: CanvasRenderingContext2D, camera: Camera) {
        
    }
}