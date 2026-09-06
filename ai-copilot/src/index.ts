import { loadConfig, getSanitizedConfig } from "./config.js";
import { createGatewayServer } from "./server.js";

function main() {
  try {
    const config = loadConfig();
    console.info("Starting AI Gateway with config:", JSON.stringify(getSanitizedConfig(config), null, 2));

    const server = createGatewayServer(config);

    server.listen(config.AI_GATEWAY_PORT, config.AI_GATEWAY_HOST, () => {
      console.info(`🚀 AI Gateway listening on http://${config.AI_GATEWAY_HOST}:${config.AI_GATEWAY_PORT}`);
    });

    const shutdown = () => {
      console.info("Shutting down AI Gateway...");
      server.close(() => {
        console.info("AI Gateway stopped cleanly.");
        process.exit(0);
      });
    };

    process.on("SIGINT", shutdown);
    process.on("SIGTERM", shutdown);
  } catch (err: any) {
    console.error("Fatal AI Gateway startup error:", err.message);
    process.exit(1);
  }
}

main();
