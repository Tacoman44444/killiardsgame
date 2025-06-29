import { PositionComponent } from "./Components";

export class Camera {
    follow: PositionComponent;
    width: number;
    height: number;
    
    constructor (follow: PositionComponent, width: number, height: number) {
        this.follow = follow;
        this.width = width;
        this.height = height;
    }

    SwitchFollow(newFollow: PositionComponent) {
        this.follow = newFollow;
    }
}