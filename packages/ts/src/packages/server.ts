// TypeScript Server Implementation

import { XProtocolProxyChannel, XProtocolProxyChannelType } from "./proxy";
import * as http from "http";

export type XProtocolServerCallType = {
  name: string;
  handler: (payload: any) => Promise<XProtocolCallResponseType>;
};

export type XProtocolCallRequestType = {
  proxy_channel_name?: string;
  name: string;
  payload: any;
  from_proxy_channel?: XProtocolProxyChannelType;
};

export type XProtocolCallResponseType = {
  success: boolean;
  data: any;
  error?: string;
};

export class XProtocolServer {
  private host: string;
  private port: number;
  private calls: Map<string, XProtocolServerCallType>;
  private proxyChannels: Map<string, XProtocolProxyChannelType>;

  constructor(host: string, port: number) {
    this.host = host;
    this.port = port;
    this.calls = new Map<string, XProtocolServerCallType>();
    this.proxyChannels = new Map<string, XProtocolProxyChannelType>();
  }

  public start() {
    http
      .createServer((req, res) => this.handleRequest(req, res))
      .listen(this.port, this.host, () => {
        console.log(`Server is running on ${this.host}:${this.port}`);
      });
  }

  public handleRequest(req: http.IncomingMessage, res: http.ServerResponse) {
    if (req.url === "/calls" && req.method === "POST") {
      let body = "";

      req.on("data", (chunk) => {
        body += chunk.toString();
      });

      let intervalId: NodeJS.Timeout = setInterval(() => {
        res.writeHead(404);
        res.end();
      }, 10000);

      req.on("end", async () => {
        try {
          const data: XProtocolCallRequestType = JSON.parse(body);

          // eğer proxy channel ise, proxy channel'a istek gönder
          if (data.proxy_channel_name && data.proxy_channel_name !== "") {
            const proxyChannel = this.proxyChannels.get(
              data.proxy_channel_name
            );
            if (proxyChannel) {
              const response = await proxyChannel.call(data.name, {
                name: data.name,
                payload: data.payload,
                from_proxy_channel: data.from_proxy_channel,
              } as XProtocolCallRequestType);
              if (response.error) {
                if (response.proxyServerError) {
                  res.writeHead(response.proxyStatus!, {
                    "Content-Type": "plain/text",
                  });
                  res.write(response.proxyError);
                  res.end();
                  return;
                } else {
                  res.writeHead(500, {
                    "Content-Type": "plain/text",
                  });
                  res.end(response.error);
                  return;
                }
              } else {
                res.writeHead(200, { "Content-Type": "application/json" });
                res.end(
                  JSON.stringify({
                    success: true,
                    data: response.data,
                    error: null,
                  })
                );
                return;
              }
            } else {
              res.writeHead(404);
              res.write("Proxy channel not found");
              res.end();
              return;
            }
          }

          // eğer call varsa, call'ı çağır
          else if (this.calls.has(data.name)) {
            if (process.env.APP_MODE === "development") {
              console.log(
                "Call isteği alındı -> " +
                  data.name +
                  " | from proxy: " +
                  data.from_proxy_channel
              );
            }

            const call = this.calls.get(data.name);
            if (call) {
              const response = await call.handler(data.payload);
              res.writeHead(200, { "Content-Type": "application/json" });
              res.end(JSON.stringify(response));
              return;
            }
          }
          // eğer call yoksa, 404 dön
          else {
            res.writeHead(404);
            res.write("Call not found");
            res.end();
            return;
          }
        } catch (err) {
          console.error(err);

          res.writeHead(400, { "Content-Type": "application/json" });
          res.end(
            JSON.stringify({ success: false, data: {}, error: "Invalid JSON" })
          );
        } finally {
          clearInterval(intervalId);
        }
      });
    } else {
      res.writeHead(404);
      res.end();
    }
  }

  public registerCall<T>(
    name: string,
    handler: (payload: T) => Promise<XProtocolCallResponseType>
  ) {
    this.calls.set(name, { name, handler });
    console.log(`Call registered -> ${name} ✔️`);
  }

  public registerProxyChannel(name: string, host: string, port: number) {
    this.proxyChannels.set(name, new XProtocolProxyChannel(name, host, port));
    console.log(`Proxy channel registered -> ${name} | ${host}:${port} ✔️`);
  }
}

export const NewXProtocolServer = (host: string, port: number) => {
  return new XProtocolServer(host, port);
};
