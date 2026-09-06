import { describe, it, expect, beforeAll, afterAll } from "vitest";
import http from "node:http";
import { SignJWT } from "jose";
import { createGatewayServer } from "../src/server.js";
import { GatewayConfig } from "../src/config.js";
import { ALLOWED_MCP_TOOLS } from "../src/mcp.js";

describe("AI Gateway Server & Endpoints", () => {
  const secret = "test-secret-key-abcdefghijklmnopqrstuvwxyz";
  const config: GatewayConfig = {
    AI_GATEWAY_PORT: 9988,
    AI_GATEWAY_HOST: "127.0.0.1",
    JWT_SECRET: secret,
    AI_MODEL_BASE_URL: "https://api.openai.com/v1",
    AI_MODEL_API_KEY: "sk-mock",
    AI_MODEL_ID: "gpt-4o-mini",
    MCP_STREAMABLE_HTTP_URL: "http://127.0.0.1:8009/mcp",
  };

  let server: http.Server;

  beforeAll(async () => {
    server = createGatewayServer(config);
    await new Promise<void>((resolve) => server.listen(9988, "127.0.0.1", () => resolve()));
  });

  afterAll(async () => {
    await new Promise<void>((resolve) => server.close(() => resolve()));
  });

  it("enforces strict 3-tool whitelist", () => {
    expect(ALLOWED_MCP_TOOLS.size).toBe(3);
    expect(ALLOWED_MCP_TOOLS.has("iot_query_devices")).toBe(true);
    expect(ALLOWED_MCP_TOOLS.has("iot_get_device_detail")).toBe(true);
    expect(ALLOWED_MCP_TOOLS.has("iot_query_telemetry_history")).toBe(true);
    expect(ALLOWED_MCP_TOOLS.has("iot_send_command")).toBe(false);
    expect(ALLOWED_MCP_TOOLS.has("iot_delete_device")).toBe(false);
  });

  it("responds to /healthz and /readyz", async () => {
    const healthRes = await fetch("http://127.0.0.1:9988/healthz");
    expect(healthRes.status).toBe(200);
    const healthJson = (await healthRes.json()) as any;
    expect(healthJson.status).toBe("ok");

    const readyRes = await fetch("http://127.0.0.1:9988/readyz");
    expect(readyRes.status).toBe(200);
    const readyJson = (await readyRes.json()) as any;
    expect(readyJson.status).toBe("ready");
  });

  it("rejects unauthenticated POST /v1/ai/chat with 401", async () => {
    const res = await fetch("http://127.0.0.1:9988/v1/ai/chat", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ messages: [{ role: "user", content: "hello" }] }),
    });
    expect(res.status).toBe(401);
  });

  it("rejects chat request with empty messages array", async () => {
    const key = new TextEncoder().encode(secret);
    const token = await new SignJWT({ UserId: "u1", OrganizationID: 1 })
      .setProtectedHeader({ alg: "HS256" })
      .setExpirationTime("1h")
      .sign(key);

    const res = await fetch("http://127.0.0.1:9988/v1/ai/chat", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ messages: [] }),
    });
    expect(res.status).toBe(400);
  });
});
