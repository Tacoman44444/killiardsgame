//maps will be created by the server.
//the server will return:
// a list containing binary --> 0: walkable   1. abyss
//

import { SpriteComponent } from "./Components";

export type MapGenData = {
    arena: number[][],
}

const pixelSize = 16;

export class Arena {
    size: number;
    mapData: MapGenData;
    walkableSprite: SpriteComponent;
    abyssSprite: SpriteComponent;

    constructor(walkableSprite: SpriteComponent, abyssSprite: SpriteComponent) {
        console.log("running map constructor");
        this.mapData = placeholderArenaGenerator();
        this.walkableSprite = walkableSprite;
        this.abyssSprite = abyssSprite;
        this.size = 2;
    }

    render(centerX: number, centerY: number, ctx: CanvasRenderingContext2D) {
        for (let row = 0; row < this.size; row++) {
            for (let column = 0; column < this.size; column++) {
                let x = centerX - (((this.size / 2)- row) * pixelSize);
                let y = centerY - (((this.size / 2) - column) * pixelSize);
                if (this.mapData.arena[row][column] == 0) {
                    ctx.drawImage(this.walkableSprite.img, this.walkableSprite.sprite.xOffset, this.walkableSprite.sprite.yOffset, this.walkableSprite.sprite.width, this.walkableSprite.sprite.height, x, y, pixelSize, pixelSize)
                } else if (this.mapData.arena[row][column] == 1) {
                    ctx.drawImage(this.abyssSprite.img, this.abyssSprite.sprite.xOffset, this.abyssSprite.sprite.yOffset, this.abyssSprite.sprite.width, this.walkableSprite.sprite.height, x, y, pixelSize, pixelSize)
                }
            }
        }
    }

}

function placeholderArenaGenerator() : MapGenData {
    return {arena: [[0, 0], [1, 0]]};
}