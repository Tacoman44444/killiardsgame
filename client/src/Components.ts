import {images} from "./AssetLoader.js"

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
    x: number;
    y: number;
    constructor(x: number, y: number) {
        this.x = x;
        this.y = y;
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