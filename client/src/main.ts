import { loadAll } from "./AssetLoader.js"
import { SpriteComponent } from "./Components.js";
import { Arena } from "./Arena.js";

async function main() {

    await loadAll();

    const canvas = document.getElementById('gameCanvas') as HTMLCanvasElement;
    const ctx = canvas.getContext('2d')!;
    console.log("works");
    const temp = new SpriteComponent("floorTile");
    const black = new SpriteComponent("abyssTile");
    //const arena = new Arena(temp, black);
    //arena.render(400, 300, ctx);
  }

main().catch(err => {
  console.error(err);
  alert("Failed to start the game. See console for details.");
});