export interface UserContext {
  rawToken: string;
}

export function extractAuthToken(
  authHeader: string | null | undefined
): UserContext {
  if (!authHeader || !authHeader.startsWith("Bearer ")) {
    throw new Error("Missing or invalid Authorization header");
  }

  const token = authHeader.slice(7).trim();
  if (!token) {
    throw new Error("Bearer token is empty");
  }

  return {
    rawToken: token,
  };
}
