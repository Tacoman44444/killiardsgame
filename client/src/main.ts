import { loadAll } from "./AssetLoader.js"
import { Game } from "./game-states.js";
import { SocketEventManager } from "./socketevent-manager.js";
import { SocketManager } from "./socket-manager.js";

async function main() {

    await loadAll();
    let wsEventManager = new SocketEventManager();
    let ws = new SocketManager("ws://localhost:8000/ws", wsEventManager);
    await ws.connect();
    const game = new Game(wsEventManager);
    game.setupInputListeners();
    console.log("time to play the game");
    game.Play();
    
  }

main().catch(err => {
  console.error(err);
  alert("Failed to start the game. See console for details.");
});