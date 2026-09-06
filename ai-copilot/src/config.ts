import { z } from "zod";
import dotenv from "dotenv";

dotenv.config();

const configSchema = z.object({
  AI_GATEWAY_PORT: z.coerce.number().default(8005),
  AI_GATEWAY_HOST: z.string().default("0.0.0.0"),
  JWT_SECRET: z.string().min(1, "JWT_SECRET is required"),
  AI_MODEL_BASE_URL: z.string().url("AI_MODEL_BASE_URL must be a valid URL"),
  AI_MODEL_API_KEY: z.string().min(1, "AI_MODEL_API_KEY is required"),
  AI_MODEL_ID: z.string().min(1, "AI_MODEL_ID is required"),
  MCP_STREAMABLE_HTTP_URL: z.string().url("MCP_STREAMABLE_HTTP_URL must be a valid URL"),
});

export type GatewayConfig = z.infer<typeof configSchema>;

export function loadConfig(env: Record<string, string | undefined> = process.env): GatewayConfig {
  const result = configSchema.safeParse(env);
  if (!result.success) {
    const issues = (result.error as any).issues || (result.error as any).errors || [];
    const errors = issues.map((e: any) => `${e.path.join(".")}: ${e.message}`).join(", ");
    throw new Error(`Invalid AI Gateway configuration: ${errors || result.error.message}`);
  }
  return result.data;
}

export function getSanitizedConfig(config: GatewayConfig): Record<string, unknown> {
  return {
    AI_GATEWAY_PORT: config.AI_GATEWAY_PORT,
    AI_GATEWAY_HOST: config.AI_GATEWAY_HOST,
    AI_MODEL_BASE_URL: config.AI_MODEL_BASE_URL,
    AI_MODEL_ID: config.AI_MODEL_ID,
    MCP_STREAMABLE_HTTP_URL: config.MCP_STREAMABLE_HTTP_URL,
    JWT_SECRET: "***",
    AI_MODEL_API_KEY: "***",
  };
}
