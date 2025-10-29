import { apiRequest } from "../../lib/api-client";
import type { StoredAuth } from "../../lib/storage";
import type { AuthCredentials } from "./schemas";
import { authResponseSchema } from "./schemas";

const LOGIN_PATH = "/api/auth/login";
const REGISTER_PATH = "/api/auth/register";

function mapCredentials(credentials: AuthCredentials) {
  return {
    cluster_id: credentials.clusterId,
    password: credentials.password,
  };
}

function mapResponse(dto: unknown): StoredAuth {
  const parsed = authResponseSchema.parse(dto);
  return {
    clusterId: parsed.cluster_id,
    token: parsed.token,
  };
}

export async function login(credentials: AuthCredentials): Promise<StoredAuth> {
  const dto = await apiRequest(LOGIN_PATH, {
    method: "POST",
    body: mapCredentials(credentials),
  });
  return mapResponse(dto);
}

export async function register(credentials: AuthCredentials): Promise<StoredAuth> {
  const dto = await apiRequest(REGISTER_PATH, {
    method: "POST",
    body: mapCredentials(credentials),
  });
  return mapResponse(dto);
}
