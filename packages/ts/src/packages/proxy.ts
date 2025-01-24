// TypeScript Proxy Implementation

import { XProtocolCallRequestType } from "./server";

export type XProtocolProxyChannelType = {
  name: string;
  host: string;
  port: number;
  call: (
    name: string,
    xprotoCallRequest: XProtocolCallRequestType
  ) => Promise<XProtocolProxyCallResponseType>;
};

export type XProtocolProxyServiceType = {
  name: string;
  host: string;
  port: number;
};

export type XProtocolProxyCallResponseType = {
  success: boolean;
  data: string;
  error?: string;
  proxyStatus?: number;
  proxyError?: string;
  proxyServerError: boolean;
};

export class XProtocolProxyChannel implements XProtocolProxyChannel {
  name: string;
  host: string;
  port: number;

  constructor(name: string, host: string, port: number) {
    this.name = name;
    this.host = host;
    this.port = port;
  }

  public async call(
    name: string,
    xprotoCallRequest: XProtocolCallRequestType
  ): Promise<XProtocolProxyCallResponseType> {
    const bodyTextJsonBytes = JSON.stringify(xprotoCallRequest);

    if (process.env.APP_MODE === "development") {
      console.log(
        "Proxy isteği yönlendirdi -> " +
          `http://${this.host}:${this.port}/calls`
      );
    }

    try {
      const response = await fetch(`http://${this.host}:${this.port}/calls`, {
        method: "POST",
        body: bodyTextJsonBytes,
      });

      if (!response.ok) {
        return {
          success: false,
          data: "",
          error: response.statusText,
          proxyStatus: response.status,
          proxyError: response.statusText,
          proxyServerError: true,
        } as XProtocolProxyCallResponseType;
      }

      const responseBody =
        (await response.json()) as XProtocolProxyCallResponseType;

      return responseBody;
    } catch (error) {
      return {
        success: false,
        data: "",
        error: error instanceof Error ? error.message : error,
        proxyStatus: undefined,
        proxyError: error instanceof Error ? error.message : error,
        proxyServerError: false,
      } as XProtocolProxyCallResponseType;
    }
  }
}
