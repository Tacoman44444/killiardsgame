//all gameobject implementations

import { Camera } from "./camera.js";
import { PositionComponent, SpriteComponent, VelocityComponent } from "./Components.js";
import { Circle, Rect, Vector2 } from "./physics.js";

interface GameObject {

    name: string;
    render(ctx : CanvasRenderingContext2D, camera: Camera) : void;
}

export class Puck implements GameObject {
    
    name: string;
    id: string;
    username: string;
    puckSprite: SpriteComponent;
    circle: Circle;

    constructor(id: string, username: string, posX: number, posY: number) {
        this.name = "puck";
        this.id = id;
        this.username = username;
        this.puckSprite = new SpriteComponent("puck");
        this.circle = new Circle(new Vector2(posX, posY), 16)
    }

    render(ctx: CanvasRenderingContext2D, camera: Camera) {

        //console.log("circle center:", this.circle.center);
        //console.log("camera follow:", camera.follow);
        //console.log("sprite dimensions:", this.puckSprite.sprite.width, this.puckSprite.sprite.height);
        const x = this.circle.center.x - camera.follow.x + (camera.width / 2) - (this.puckSprite.sprite.width / 2);
        const y = this.circle.center.y - camera.follow.y + (camera.height / 2) - (this.puckSprite.sprite.height / 2);
        //console.log("x coord of puck: ", x)
        //console.log("y coord of puck: ", y)
        ctx.drawImage(
            this.puckSprite.img,
            this.puckSprite.sprite.xOffset, this.puckSprite.sprite.yOffset,
            this.puckSprite.sprite.width, this.puckSprite.sprite.height,
            x, y,
            32, 32
        );
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

    render(ctx: CanvasRenderingContext2D, camera: Camera) {

    }
}

export class Wall implements GameObject {

    name: string;
    wallSprite: SpriteComponent;
    rect: Rect;

    constructor(posX: number, posY: number) {
        this.name = "wall";
        this.wallSprite = new SpriteComponent("wallTile");
        this.rect = new Rect(new Vector2(posX, posY), 64, 64);
    }

    render(ctx: CanvasRenderingContext2D, camera: Camera) {
        const x = this.rect.topLeft.x - camera.follow.x + (camera.width / 2);
        const y = this.rect.topLeft.y - camera.follow.y + (camera.height / 2);

        ctx.drawImage(
            this.wallSprite.img,
            this.wallSprite.sprite.xOffset, this.wallSprite.sprite.yOffset,
            this.wallSprite.sprite.width, this.wallSprite.sprite.height,
            x, y,
            64, 64
        );
    }
}