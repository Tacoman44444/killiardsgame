import { Board } from "./ui"
import { World } from "./World"


interface GameState {

    name:       string
    board:      Board
    world:      World
    processInput() : void;
    update() : void;
    render(ctx: CanvasRenderingContext2D) : void;
    Enter() : void;
}