import { jwtVerify } from "jose";

export interface UserContext {
  userId: string;
  organizationId: number;
  rawToken: string;
}

export async function verifyAuthToken(
  authHeader: string | null | undefined,
  jwtSecret: string
): Promise<UserContext> {
  if (!authHeader || !authHeader.startsWith("Bearer ")) {
    throw new Error("Missing or invalid Authorization header");
  }

  const token = authHeader.slice(7).trim();
  if (!token) {
    throw new Error("Bearer token is empty");
  }

  const secretKey = new TextEncoder().encode(jwtSecret);
  const { payload } = await jwtVerify(token, secretKey);

  const userId =
    (payload.UserId as string) ||
    (payload.userId as string) ||
    (payload.user_id as string) ||
    (payload.sub as string) ||
    "";

  const orgId =
    (payload.OrganizationID as number) ||
    (payload.organizationId as number) ||
    (payload.organization_id as number) ||
    0;

  if (!orgId || orgId <= 0) {
    throw new Error("Token missing valid organization identity");
  }

  return {
    userId,
    organizationId: Number(orgId),
    rawToken: token,
  };
}
