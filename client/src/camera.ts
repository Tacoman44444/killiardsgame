import { PositionComponent } from "./Components";
import { Vector2 } from "./physics";

export class Camera {
    follow: Vector2;
    width: number;
    height: number;
    
    constructor (follow: Vector2, width: number, height: number) {
        this.follow = follow;
        this.width = width;
        this.height = height;
    }

    SwitchFollow(newFollow: Vector2) {
        this.follow = newFollow;
    }
}