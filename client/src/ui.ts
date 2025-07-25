import { SpriteComponent } from "./Components.js";
import { GameInput } from "./game-states.js";

export class Board {

    buttons: Button[] = [];
    textButtons: TextButton[] = [];
    inputBoxes: InputTextBox[] = [];
    displayBoxes: DisplayTextBox[] = []
    images: Img[] = []

    addButton(name: string, x: number, y: number, width: number, height: number, sprite: SpriteComponent, onClick: () => void, visible: boolean) {
        const btn = new Button(name, x, y, width, height, sprite, onClick, visible)
        this.buttons.push(btn)
    }

    addTextButton(name: string, x: number, y: number, width: number, height: number, text: string, onClick: () => void, visible: boolean, fontSize: number) {
        const btn = new TextButton(name, x, y, width, height, text, onClick, visible, fontSize)
        this.textButtons.push(btn)
    }

    addImage(name: string, x: number, y: number, width: number, height: number, sprite: SpriteComponent, visible: boolean) {
        const img = new Img(name, x, y, width, height, sprite, visible)
        this.images.push(img)
    }

    addDisplayTextBox(name: string, x: number, y: number, width: number, height: number, text: string, fontSize: number, visible: boolean) {
        const box = new DisplayTextBox(name, x, y, width, height, text, fontSize, visible);
        this.displayBoxes.push(box)
    }

    addInputTextBox(name: string, x: number, y: number, width: number, height: number, maxLength: number) {
        const box = new InputTextBox(name, x, y, width, height, maxLength);
        this.inputBoxes.push(box);
    }


    processInput(input: any) {
        this.textButtons.forEach((btn) => btn.processInput(input));
        this.buttons.forEach((btn) => btn.processInput(input));
        this.inputBoxes.forEach((box) => box.processInput(input));
    }

    render(ctx: CanvasRenderingContext2D) : void {
        this.buttons.forEach((btn) => btn.render(ctx));
        this.textButtons.forEach((btn) => btn.render(ctx))
        this.inputBoxes.forEach((box) => box.render(ctx));
        this.displayBoxes.forEach((box) => box.render(ctx));
        this.images.forEach((img) => img.render(ctx))
    }

    updateDisplayTextBox(boxName: string, newText: string) {
        this.displayBoxes.forEach((box) => {
            if (box.name == boxName) {
                box.UpdateText(newText)
            }
        })
    }

    getDisplayTextBox(id: string): DisplayTextBox | undefined {
        for (const box of this.displayBoxes) {
            if (box.name === id) {
                return box;
            }
        }
    }

    getTextButton(id: string): TextButton | undefined {
        for (const btn of this.textButtons) {
            if (btn.name == id) {
                return btn;
            }
        }
    }

}

export class BoardEventManager {
    board: Board;
    eventMap: { [key: string]: Function } = {};

    constructor(board: Board) {
        this.board = board;
    }

    addEvent(eventType: string, callback: Function) {
        this.eventMap[eventType] = callback;
    }

    onEvent(eventType: string, data?: any) {
        console.log("=== BoardEventManager.onEvent ===");
        console.log("eventType:", eventType);
        console.log("eventMap keys:", Object.keys(this.eventMap));
        console.log("eventMap[eventType]:", this.eventMap[eventType]);
        this.eventMap[eventType](data);
    }
}

type uiElement = Button | TextButton | InputTextBox | DisplayTextBox;

export class Img {
    name: string;
    x: number;
    y: number;
    width: number; 
    height: number;
    sprite: SpriteComponent;
    visible: boolean;

    constructor(name: string, x: number, y: number, width: number, height: number, sprite: SpriteComponent, visible: boolean) {
        this.name = name;
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.sprite = sprite;
        this.visible = visible
    }

    render(ctx: CanvasRenderingContext2D) {
        if (this.visible) {
            ctx.drawImage(this.sprite.img, this.sprite.sprite.xOffset, this.sprite.sprite.yOffset, this.sprite.sprite.width, this.sprite.sprite.height, this.x, this.y, this.width, this.height)
        }
    }
}

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
        this.visible = visible;
    }

    processInput(input: GameInput) {
        if (this.visible) {
            if (input.type == "mouseup") {
                if (contains(this, input.cameraX, input.cameraY)) {
                    this.onClick();
                }
            }
        }

    }

    render(ctx: CanvasRenderingContext2D) {
        if (this.visible) {
            ctx.drawImage(this.sprite.img, this.sprite.sprite.xOffset, this.sprite.sprite.yOffset, this.sprite.sprite.width, this.sprite.sprite.height, this.x, this.y, this.width, this.height)
        }
    }
};

class TextButton {
    name: string;
    x: number;
    y: number;
    width: number;
    height: number;
    text: string;
    onClick: () => void;
    visible: boolean;
    fontSize: number;
    fontFamily: string;

    constructor(
        name: string,
        x: number,
        y: number,
        width: number,
        height: number,
        text: string,
        onClick: () => void,
        visible: boolean,
        fontSize: number = 24,
        fontFamily: string = "'Orbitron'"
    ) {
        this.name = name;
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.text = text;
        this.onClick = onClick;
        this.visible = visible;
        this.fontSize = fontSize;
        this.fontFamily = fontFamily;
    }

    processInput(input: GameInput) {
        if (this.visible && input.type === "mouseup") {
            if (contains(this, input.cameraX, input.cameraY)) {
                this.onClick();
            }
        }
    }

    render(ctx: CanvasRenderingContext2D) {
        if (!this.visible) return;

        // Background
        ctx.fillStyle = "#333";
        ctx.fillRect(this.x, this.y, this.width, this.height);

        // Border
        ctx.strokeStyle = "#999";
        ctx.lineWidth = 2;
        ctx.strokeRect(this.x, this.y, this.width, this.height);

        // Text
        ctx.fillStyle = "white";
        ctx.font = `${this.fontSize}px ${this.fontFamily}`;
        ctx.textAlign = "center";
        ctx.textBaseline = "middle";
        ctx.fillText(this.text, this.x + this.width / 2, this.y + this.height / 2);
    }
}

class DisplayTextBox {
    name: string;
    x: number;
    y: number;
    width: number;
    height: number;
    text: string = "-";
    fontSize: number;
    visible: boolean;

    constructor(name: string, x: number, y: number, width: number, height: number, text: string, fontSize: number = 16, visible: boolean) {
        this.name = name;
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.text = text;
        this.fontSize = fontSize;
        this.visible = visible
    }

    UpdateText(newText: string) {
        this.text = newText;
    }

    render(ctx: CanvasRenderingContext2D) {

        
        // Set up text styling
        if (this.visible) {

        ctx.fillStyle = "#999999"; // text color
        ctx.font = `${this.fontSize}px 'Orbitron'`;   // dynamic font size
        ctx.textBaseline = "top";

        // Word-wrap logic
        const padding = 0;
        const maxWidth = this.width - 2 * padding;
        const words = (this.text ?? "").split(" ");
        let line = "";
        const lineHeight = this.fontSize + 4;
        let y = this.y + padding;

        for (let i = 0; i < words.length; i++) {
            const testLine = line + words[i] + " ";
            const testWidth = ctx.measureText(testLine).width;

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
            this.focused = contains(this, input.cameraX, input.cameraY);
        } else if (input.type == "keydown" && this.focused) {
            const key = input.event.key;

            if (key === "Backspace") {
                this.text = this.text.slice(0, -1); // remove last character
            } else if (isAlphanumeric(key) && this.text.length < this.maxLength) {
                this.text += key;
            }
        }
    }

    render(ctx: CanvasRenderingContext2D) {
        const padding = 10;
        const centerX = this.x + this.width / 2;
        const centerY = this.y + this.height / 2;

        // Draw background
        ctx.fillStyle = this.focused ? "#222" : "#444";
        ctx.fillRect(this.x, this.y, this.width, this.height);

        // Draw border
        ctx.strokeStyle = this.focused ? "#00f" : "#888";
        ctx.lineWidth = 2;
        ctx.strokeRect(this.x, this.y, this.width, this.height);

        // Set font
        ctx.fillStyle = "white";
        ctx.font = "24px 'Orbitron'";
        ctx.textBaseline = "middle";
        ctx.textAlign = "left";

        // Measure and clip text if needed
        let textToRender = this.text;
        let measuredWidth = ctx.measureText(textToRender).width;

        while (measuredWidth > this.width - 2 * padding && textToRender.length > 0) {
            textToRender = textToRender.slice(1); // trim from front
            measuredWidth = ctx.measureText(textToRender).width;
        }

        const textX = this.x + padding;
        const textY = centerY;

        // Draw the text
        ctx.fillText(textToRender, textX, textY);

        // Draw caret if focused
        if (this.focused && Date.now() % 1000 < 500) {
            const caretX = textX + measuredWidth;
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