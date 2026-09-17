import { describe, it, expect } from "vitest";
import { extractAuthToken } from "../src/auth.js";

describe("Gateway Auth Extraction", () => {
  it("successfully extracts valid Bearer token", () => {
    const rawToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.test-token-payload";
    const context = extractAuthToken(`Bearer ${rawToken}`);
    expect(context.rawToken).toBe(rawToken);
  });

  it("fails on missing Authorization header", () => {
    expect(() => extractAuthToken(undefined)).toThrow(
      "Missing or invalid Authorization header"
    );
  });

  it("fails on non-Bearer Authorization header", () => {
    expect(() => extractAuthToken("Basic dXNlcjpwYXNz")).toThrow(
      "Missing or invalid Authorization header"
    );
  });

  it("fails on empty Bearer token", () => {
    expect(() => extractAuthToken("Bearer   ")).toThrow(
      "Bearer token is empty"
    );
  });
});
