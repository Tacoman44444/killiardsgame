import {images} from "./AssetLoader.js"
import { Vector2 } from "./physics.js";

export interface SpriteDef {
    src: string;
    xOffset: number;
    yOffset: number;
    width: number;
    height: number;
}

const SPRITE_MAP: Record <string, SpriteDef> = {
    floorTile : { src: "block_sheet", xOffset: 256, yOffset: 0, width: 16, height: 16 },
    abyssTile : { src: "block_sheet", xOffset: 272, yOffset: 0, width: 16, height: 16 },
    wallTile  : { src: "block_sheet", xOffset: 0, yOffset: 0, width: 16, height: 16 },
    brickTile : { src: "block_sheet", xOffset: 288, yOffset: 0, width: 16, height: 16 },
    puck      : { src: "puck", xOffset: 0, yOffset: 0, width: 128, height: 128 },
    createroom_button: { src: "createroom_button", xOffset: 0, yOffset: 0, width: 68, height: 16 },
}

export class SpriteComponent {
    readonly sprite: SpriteDef;
    readonly img: HTMLImageElement;
    constructor(id: string) {
        this.sprite = SPRITE_MAP[id]
        const img = images.get(this.sprite.src);
        if (!img) throw new Error(`Sprite '${id}' requires missing image '${this.sprite.src}'`);
        this.img = img;
    }

}

export class PositionComponent {
    pos: Vector2;
    constructor(x: number, y: number) {
        this.pos = new Vector2(x, y)
    }
}

export class VelocityComponent {
    x: number;
    y: number;
    constructor(x: number, y: number) {
        this.x = x;
        this.y = y;
    }
}