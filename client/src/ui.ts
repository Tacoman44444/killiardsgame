import { SpriteComponent } from "./Components";


export class Board {

    buttons: Button[] = [];
    textBoxes: TextBox[] = [];

    addButton(name: string, x: number, y: number, width: number, height: number, sprite: SpriteComponent, onClick: () => void, visible: boolean) {
        const btn = new Button(name, x, y, width, height, sprite, onClick, visible)
        this.buttons.push(btn)
    }

    addTextBox(name: string, x: number, y: number, width: number, height: number, maxLength: number) {
        const box = new TextBox(name, x, y, width, height, maxLength);
        this.textBoxes.push(box);
    }

    processInput(input: any) {
        
    }

    render(ctx: CanvasRenderingContext2D) : void {
        this.buttons.forEach((btn) => btn.render(ctx));
        this.textBoxes.forEach((box) => box.render(ctx));
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
        this.visible = true;
    }

    render(ctx: CanvasRenderingContext2D) {
        ctx.drawImage(this.sprite.img, this.sprite.sprite.xOffset, this.sprite.sprite.yOffset, this.sprite.sprite.width, this.sprite.sprite.height, this.x, this.y, this.width, this.height)
    }
};

class TextBox {
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