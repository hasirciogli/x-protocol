// TypeScript Client Implementation

export interface XProtocolClientCallResponse {
  success: boolean;
  data: any; // json.RawMessage
  error?: string;
}

export interface XProtocolClientCallRequest {
  proxyChannelName?: string; // proxy channel name for proxy mode xprotocol server
  name: string; // function name
  payload: any; // json.RawMessage
}

export class XProtocolClient {
  private host: string;
  private port: number;

  constructor(host: string, port: number) {
    this.host = host;
    this.port = port;
  }

  public async call(
    request: XProtocolClientCallRequest
  ): Promise<XProtocolClientCallResponse> {
    const bodyTextJson = JSON.stringify(request);

    if (process.env.APP_MODE === "development") {
      console.log(`Call isteği gönderildi -> ${request.name}`);
    }

    try {
      const response = await fetch(`http://${this.host}:${this.port}/calls`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: bodyTextJson,
      });

      const responseBody = await response.json();

      if (!response.ok) {
        return {
          success: false,
          data: null,
          error: responseBody.error || "Unknown error",
        };
      }

      if (process.env.APP_MODE === "development") {
        console.log(
          `Call yanıtı alındı -> ${JSON.stringify(responseBody.data)}`
        );
      }

      return {
        success: true,
        data: responseBody.data,
        error: responseBody.error,
      };
    } catch (error) {
      return {
        success: false,
        data: null,
        error: error instanceof Error ? error.message : "Unknown error",
      };
    }
  }
}

export function newXProtocolClient(
  host: string,
  port: number
): XProtocolClient {
  return new XProtocolClient(host, port);
}
