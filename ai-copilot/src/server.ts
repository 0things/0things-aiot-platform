import http from "node:http";
import { GatewayConfig } from "./config.js";
import { verifyAuthToken } from "./auth.js";
import { handleChatRequest } from "./chat.js";

function setCorsHeaders(res: http.ServerResponse) {
  res.setHeader("Access-Control-Allow-Origin", "*");
  res.setHeader("Access-Control-Allow-Methods", "GET, POST, OPTIONS");
  res.setHeader("Access-Control-Allow-Headers", "Authorization, Content-Type");
}

export function createGatewayServer(config: GatewayConfig) {
  const server = http.createServer(async (req, res) => {
    const startTime = Date.now();
    const requestId = crypto.randomUUID();
    const url = new URL(req.url || "/", `http://${req.headers.host || "localhost"}`);
    setCorsHeaders(res);

    if (req.method === "OPTIONS") {
      res.writeHead(204);
      res.end();
      return;
    }

    // Health & Readiness checks
    if (req.method === "GET" && (url.pathname === "/healthz" || url.pathname === "/health")) {
      res.writeHead(200, { "Content-Type": "application/json" });
      res.end(JSON.stringify({ status: "ok", timestamp: Date.now() }));
      return;
    }

    if (req.method === "GET" && (url.pathname === "/readyz" || url.pathname === "/ready")) {
      res.writeHead(200, { "Content-Type": "application/json" });
      res.end(JSON.stringify({ status: "ready" }));
      return;
    }

    // Main AI Chat Route
    if (req.method === "POST" && url.pathname === "/v1/ai/chat") {
      let userContext;
      try {
        userContext = await verifyAuthToken(req.headers.authorization, config.JWT_SECRET);
      } catch (authErr: any) {
        console.warn(`[${requestId}] Auth failed: ${authErr.message}`);
        res.writeHead(401, { "Content-Type": "application/json" });
        res.end(JSON.stringify({ error: "Unauthorized: " + authErr.message }));
        return;
      }

      console.info(
        `[${requestId}] AI chat request started: user=${userContext.userId}, org=${userContext.organizationId}`
      );

      // Convert Node IncomingMessage to Web Standard Request
      const chunks: Buffer[] = [];
      for await (const chunk of req) {
        chunks.push(typeof chunk === "string" ? Buffer.from(chunk) : chunk);
      }
      const bodyBuffer = Buffer.concat(chunks);

      const webReq = new Request(url.toString(), {
        method: "POST",
        headers: req.headers as HeadersInit,
        body: bodyBuffer.length > 0 ? bodyBuffer : undefined,
      });

      try {
        const webRes = await handleChatRequest(webReq, userContext, config);
        res.writeHead(webRes.status, Object.fromEntries(webRes.headers.entries()));

        if (webRes.body) {
          const reader = webRes.body.getReader();
          while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            res.write(value);
          }
        }
        res.end();
        console.info(
          `[${requestId}] AI chat request finished: duration=${Date.now() - startTime}ms, status=${webRes.status}`
        );
      } catch (err: any) {
        console.error(`[${requestId}] Unhandled chat error:`, err);
        if (!res.headersSent) {
          res.writeHead(500, { "Content-Type": "application/json" });
          res.end(
            JSON.stringify({ error: "Internal server error. Please try again later." })
          );
        }
      }
      return;
    }

    res.writeHead(404, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ error: "Not Found" }));
  });

  return server;
}
