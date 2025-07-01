export const ASSET_MAP: Record<string, string> = {
    block_sheet: "assets/spritesheet.png",
    puck: "assets/disc2.png",
    createroom_button: "assets/createroom.png",
}

export const images = new Map<string, HTMLImageElement>();

export async function loadAll(): Promise<void> {

    const loadImage = (src: string): Promise<HTMLImageElement> =>
        new Promise((resolve, reject) => {
            const img = new Image();
            img.onload = () => resolve(img);
            img.onerror = reject;
            img.src = src;
        });

    const entries = Object.entries(ASSET_MAP);

    const promisies = entries.map(async ([key, path]) => {
        const image = await loadImage(path);
        images.set(key, image);
    });

    await Promise.all(promisies)
}
