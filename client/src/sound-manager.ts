export class SoundManager {
    sounds: Record<string, HTMLAudioElement> = {}

    loadSound(name: string, src: string, volume: number = 1.0) {
        const audio = new Audio(src)
        audio.volume = volume
        this.sounds[name] = audio
    }

    play(name: string) {
        const sound = this.sounds[name];
        if (sound) {
            const clone = sound.cloneNode(true) as HTMLAudioElement;
            clone.play().catch(console.error)
        }
    }
    constructor() {
        this.loadSound("dead", "assets/sounds/dead.mp3")
        this.loadSound("menuclick", "assets/sounds/menuclick.wav")
        this.loadSound("shoot", "assets/sounds/shoot.wav", 0.4)
        this.loadSound("turnactive", "assets/sounds/turnactive.wav")
    }
}