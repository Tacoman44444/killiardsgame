import { loadAll } from "./AssetLoader.js"
import { SpriteComponent } from "./Components/SpriteComponent.js";

async function main() {

    await loadAll();

    const canvas = document.getElementById('gameCanvas') as HTMLCanvasElement;
    const ctx = canvas.getContext('2d')!;
    console.log("works");
    const temp = new SpriteComponent("floorTile");
    ctx.drawImage(temp.img, temp.sprite.xOffset, temp.sprite.yOffset, temp.sprite.width, temp.sprite.height, 100, 100, temp.sprite.width, temp.sprite.height)
}

main().catch(err => {
  console.error(err);
  alert("Failed to start the game. See console for details.");
});