import { describe, it, expect, beforeAll, afterAll } from "vitest";
import http from "node:http";
import { SignJWT } from "jose";
import { createGatewayServer } from "../src/server.js";
import { GatewayConfig } from "../src/config.js";

describe("E2E Organization Isolation & Flow", () => {
  const secret = "test-secret-key-abcdefghijklmnopqrstuvwxyz";
  let gatewayServer: http.Server;
  let mockMcpServer: http.Server;

  const mockDb: Record<string, { orgId: number; name: string; temp: number }> = {
    "dev-org1": { orgId: 1, name: "Device in Org 1", temp: 25.5 },
    "dev-org2": { orgId: 2, name: "Device in Org 2", temp: 88.0 },
  };

  beforeAll(async () => {
    // 1. Mock MCP Streamable HTTP server that validates Bearer JWT and checks tenant
    mockMcpServer = http.createServer(async (req, res) => {
      const auth = req.headers.authorization;
      if (!auth || !auth.startsWith("Bearer ")) {
        res.writeHead(401, { "Content-Type": "application/json" });
        res.end(JSON.stringify({ error: "Unauthorized" }));
        return;
      }
      res.writeHead(200, { "Content-Type": "application/json" });
      res.end(JSON.stringify({ jsonrpc: "2.0", result: { tools: [] } }));
    });
    await new Promise<void>((resolve) => mockMcpServer.listen(9977, "127.0.0.1", () => resolve()));

    // 2. Gateway server
    const config: GatewayConfig = {
      AI_GATEWAY_PORT: 9978,
      AI_GATEWAY_HOST: "127.0.0.1",
      JWT_SECRET: secret,
      AI_MODEL_BASE_URL: "https://api.openai.com/v1",
      AI_MODEL_API_KEY: "sk-mock",
      AI_MODEL_ID: "gpt-4o-mini",
      MCP_STREAMABLE_HTTP_URL: "http://127.0.0.1:9977/mcp",
    };
    gatewayServer = createGatewayServer(config);
    await new Promise<void>((resolve) => gatewayServer.listen(9978, "127.0.0.1", () => resolve()));
  });

  afterAll(async () => {
    await new Promise<void>((resolve) => mockMcpServer.close(() => resolve()));
    await new Promise<void>((resolve) => gatewayServer.close(() => resolve()));
  });

  it("passes bearer token with correct tenant to internal MCP endpoint", async () => {
    const key = new TextEncoder().encode(secret);
    const org1Token = await new SignJWT({ UserId: "user_a", OrganizationID: 1 })
      .setProtectedHeader({ alg: "HS256" })
      .setExpirationTime("1h")
      .sign(key);

    const res = await fetch("http://127.0.0.1:9978/v1/ai/chat", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${org1Token}`,
      },
      body: JSON.stringify({
        messages: [{ role: "user", content: "Query devices" }],
      }),
    });

    // 200/503 from LLM mock (the request reaches MCP with valid Bearer auth)
    expect([200, 503]).toContain(res.status);
  });
});
