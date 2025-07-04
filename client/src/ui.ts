import { SpriteComponent } from "./Components.js";
import { GameInput } from "./game-states.js";


export class Board {

    buttons: Button[] = [];
    textBoxes: InputTextBox[] = [];
    displayBoxes: DisplayTextBox[] = []

    addButton(name: string, x: number, y: number, width: number, height: number, sprite: SpriteComponent, onClick: () => void, visible: boolean) {
        const btn = new Button(name, x, y, width, height, sprite, onClick, visible)
        this.buttons.push(btn)
    }

    addDisplayTextBox(name: string, x: number, y: number, width: number, height: number, text: string) {
        const box = new DisplayTextBox(name, x, y, width, height, text);
        this.displayBoxes.push(box)
    }

    addInputTextBox(name: string, x: number, y: number, width: number, height: number, maxLength: number) {
        const box = new InputTextBox(name, x, y, width, height, maxLength);
        this.textBoxes.push(box);
    }

    processInput(input: any) {
        this.buttons.forEach((btn) => btn.processInput(input));
        this.textBoxes.forEach((box) => box.processInput(input));
    }

    render(ctx: CanvasRenderingContext2D) : void {
        this.buttons.forEach((btn) => btn.render(ctx));
        this.textBoxes.forEach((box) => box.render(ctx));
    }

}

type uiElement = Button | InputTextBox;

class Button {
    name:   string;
    x:      number;
    y:      number;
    width:  number;
    height: number;
    sprite: SpriteComponent;
    onClick: () => void;
    visible: boolean; 

    constructor(name: string, x: number, y: number, width: number, height: number, sprite: SpriteComponent, onClick: () => void, visible: boolean) {
        this.name = name;
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.sprite = sprite
        this.onClick = onClick;
        this.visible = true;
    }

    processInput(input: GameInput) {
        if (input.type == "mouseup") {
            if (contains(this, input.cameraX, input.cameraY)) {
                this.onClick();
            }
        }
    }

    render(ctx: CanvasRenderingContext2D) {
        ctx.drawImage(this.sprite.img, this.sprite.sprite.xOffset, this.sprite.sprite.yOffset, this.sprite.sprite.width, this.sprite.sprite.height, this.x, this.y, this.width, this.height)
    }
};

class DisplayTextBox {
    name: string;
    x: number;
    y: number;
    width: number;
    height: number;
    text: string;

    constructor(name: string, x: number, y: number, width: number, height: number, text: string) {
        this.name = name;
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.text = text;
    }

    UpdateText(newText: string) {
        this.text = newText;
    }

    render(ctx: CanvasRenderingContext2D) {
    // Optional: Draw the background box
    ctx.fillStyle = "#ffffff"; // background color
    ctx.fillRect(this.x, this.y, this.width, this.height);

    // Optional: Draw a border
    ctx.strokeStyle = "#000000"; // border color
    ctx.strokeRect(this.x, this.y, this.width, this.height);

    // Set up text styling
    ctx.fillStyle = "#000000"; // text color
    ctx.font = "16px Arial";   // font
    ctx.textBaseline = "top";

    // Word-wrap logic
    const padding = 4;
    const maxWidth = this.width - 2 * padding;
    const words = this.text.split(" ");
    let line = "";
    let lineHeight = 20;
    let y = this.y + padding;

    for (let i = 0; i < words.length; i++) {
        const testLine = line + words[i] + " ";
        const metrics = ctx.measureText(testLine);
        const testWidth = metrics.width;

        if (testWidth > maxWidth && line !== "") {
            ctx.fillText(line, this.x + padding, y);
            line = words[i] + " ";
            y += lineHeight;
        } else {
            line = testLine;
        }
    }

    ctx.fillText(line, this.x + padding, y); // draw last line
}
}

class InputTextBox {
    name:    string;
    x:       number;
    y:       number;
    width:   number;
    height:  number;
    text:    string;
    focused: boolean;
    maxLength: number;

    constructor(name: string, x: number, y: number, width: number, height: number, maxLength: number) {
        this.name = name;
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.text = "";
        this.focused = false;
        this.maxLength = maxLength;
    }

    processInput(input: any) {
        if (input.type == "mouseup") {
            if (contains(this, input.cameraX, input.cameraY)) {
                this.focused = true;
            } else {
                this.focused = false;
            }
        } else if (input.type == "keydown") {
            if (this.focused) {
                if (isAlphanumeric(input.event.key) && this.text.length < this.maxLength) {
                    this.text += input.event.key;
                }
            }
        }
    }

    render(ctx: CanvasRenderingContext2D) {
        ctx.fillStyle = this.focused ? "#222" : "#444";
        ctx.fillRect(this.x, this.y, this.width, this.height);

        ctx.strokeStyle = this.focused ? "00f" : "#888";
        ctx.lineWidth = 2;
        ctx.strokeRect(this.x, this.y, this.width, this.height);

        ctx.fillStyle = "white"
        ctx.font = "18px sans-serif";
        ctx.textBaseline = "middle";
        ctx.textAlign = "left";
        ctx.fillText(this.text, this.x + 5, this.y + this.height / 2);

        if (this.focused && Date.now() % 1000 < 500) {
            const textWidth = ctx.measureText(this.text).width;
            const caretX = this.x + 5 + textWidth;
            ctx.beginPath();
            ctx.moveTo(caretX, this.y + 5);
            ctx.lineTo(caretX, this.y + this.height - 5);
            ctx.strokeStyle = "white";
            ctx.stroke();
        }
    }
};

function contains(ui: uiElement, inputX: number, inputY: number): boolean {
    return (
        inputX >= ui.x &&
        inputX <= ui.x + ui.width &&
        inputY >= ui.y &&
        inputY <= ui.y + ui.height
    );
}

function isAlphanumeric(key: string): boolean {
    return /^[a-zA-Z0-9]$/.test(key);
}