// TypeScript Client Implementation

export class Client {
    private serverUrl: string;

    constructor(serverUrl: string) {
        this.serverUrl = serverUrl;
    }

    public connect() {
        console.log(`Connecting to server at ${this.serverUrl}...`);
        // Burada bağlantı mantığını ekleyebilirsiniz.
    }

    public sendRequest(data: any) {
        console.log(`Sending request with data: ${JSON.stringify(data)}`);
        // Burada istek gönderme mantığını ekleyebilirsiniz.
    }
} 