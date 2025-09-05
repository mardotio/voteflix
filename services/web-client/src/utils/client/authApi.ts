import { ApiFetch } from "./apiFetch";

export interface TokenResponse {
  token: string;
  expiresAt: number;
}

export const authApi = {
  create: (loginToken: string) =>
    ApiFetch.fetch<TokenResponse>({
      method: "POST",
      route: "/api/auth/token",
      headers: { Authorization: `Bearer ${loginToken}` },
    }),
};
