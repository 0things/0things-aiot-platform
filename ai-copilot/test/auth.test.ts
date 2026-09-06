import { describe, it, expect } from "vitest";
import { SignJWT } from "jose";
import { verifyAuthToken } from "../src/auth.js";

describe("Gateway JWT Auth", () => {
  const secret = "test-secret-key-abcdefghijklmnopqrstuvwxyz";

  async function generateTestToken(userId: string, orgId: number, secretKey: string, exp = "1h") {
    const key = new TextEncoder().encode(secretKey);
    return new SignJWT({
      UserId: userId,
      OrganizationID: orgId,
    })
      .setProtectedHeader({ alg: "HS256" })
      .setIssuedAt()
      .setExpirationTime(exp)
      .sign(key);
  }

  it("successfully verifies valid Bearer token", async () => {
    const token = await generateTestToken("user_123", 42, secret);
    const context = await verifyAuthToken(`Bearer ${token}`, secret);
    expect(context.userId).toBe("user_123");
    expect(context.organizationId).toBe(42);
    expect(context.rawToken).toBe(token);
  });

  it("fails on missing Authorization header", async () => {
    await expect(verifyAuthToken(undefined, secret)).rejects.toThrow("Missing or invalid Authorization header");
  });

  it("fails on invalid signature", async () => {
    const token = await generateTestToken("user_123", 42, "wrong-secret-key-1234567890123456");
    await expect(verifyAuthToken(`Bearer ${token}`, secret)).rejects.toThrow();
  });

  it("fails when organization ID is missing or zero", async () => {
    const token = await generateTestToken("user_123", 0, secret);
    await expect(verifyAuthToken(`Bearer ${token}`, secret)).rejects.toThrow("missing valid organization identity");
  });
});
