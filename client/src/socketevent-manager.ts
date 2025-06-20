export class SocketEventManager {
    private listeners: { [key: string]: Function[] } = {};

    subscribe(event: string, callback: Function): void {
        if (!this.listeners[event]) {
            this.listeners[event] = [];
        }
        this.listeners[event].push(callback);
    }

    unsubscribe(event: string, callback: Function): void {
        if (!this.listeners[event]) return;

        this.listeners[event] = this.listeners[event].filter (
            (current) => current !== callback
        );
    }

    emit(event: string, data?: any): void {
        if (!this.listeners[event]) return;

        for (const listener of this.listeners[event]) {
            listener(data);
        }
    }

}