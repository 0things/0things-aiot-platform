import { describe, it, expect } from "vitest";
import { loadConfig, getSanitizedConfig } from "../src/config.js";

describe("Gateway Config", () => {
  const validEnv = {
    AI_GATEWAY_PORT: "8005",
    AI_GATEWAY_HOST: "0.0.0.0",
    JWT_SECRET: "my-secret-key-12345",
    AI_MODEL_BASE_URL: "https://api.openai.com/v1",
    AI_MODEL_API_KEY: "sk-test-12345",
    AI_MODEL_ID: "gpt-4o-mini",
    MCP_STREAMABLE_HTTP_URL: "http://127.0.0.1:8009/mcp",
  };

  it("successfully parses valid environment", () => {
    const config = loadConfig(validEnv);
    expect(config.AI_GATEWAY_PORT).toBe(8005);
    expect(config.AI_MODEL_ID).toBe("gpt-4o-mini");
    expect(config.JWT_SECRET).toBe("my-secret-key-12345");
  });

  it("throws error when required variables are missing", () => {
    expect(() => loadConfig({})).toThrow("Invalid AI Gateway configuration");
  });

  it("masks sensitive secrets when sanitized", () => {
    const config = loadConfig(validEnv);
    const sanitized = getSanitizedConfig(config);
    expect(sanitized.JWT_SECRET).toBe("***");
    expect(sanitized.AI_MODEL_API_KEY).toBe("***");
    expect(sanitized.AI_MODEL_ID).toBe("gpt-4o-mini");
  });
});
