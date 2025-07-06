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
    floorTile :  { src: "texture_stone", xOffset: 0, yOffset: 0, width: 1024, height: 1024 },
    decayedTile: { src: "texture_stone_dark", xOffset: 0, yOffset: 0, width: 1024, height: 1024 },
    abyssTile :  { src: "abyss_texture", xOffset: 0, yOffset: 0, width: 16, height: 16 },
    wallTile  :  { src: "block_sheet", xOffset: 0, yOffset: 0, width: 16, height: 16 },
    brickTile :  { src: "block_sheet", xOffset: 288, yOffset: 0, width: 16, height: 16 },
    puck      :  { src: "puck", xOffset: 0, yOffset: 0, width: 128, height: 128 },
    createroom_button: { src: "createroom_button", xOffset: 0, yOffset: 0, width: 68, height: 16 },
    power_level_5: { src: "powerlevel_sheet", xOffset: 0, yOffset: 0, width: 10, height: 20 },
    power_level_4: { src: "powerlevel_sheet", xOffset: 10, yOffset: 0, width: 10, height: 20 },
    power_level_3: { src: "powerlevel_sheet", xOffset: 20, yOffset: 0, width: 10, height: 20 },
    power_level_2: { src: "powerlevel_sheet", xOffset: 30, yOffset: 0, width: 10, height: 20 },
    power_level_1: { src: "powerlevel_sheet", xOffset: 40, yOffset: 0, width: 10, height: 20 },
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