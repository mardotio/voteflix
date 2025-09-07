import { ApiConfig } from "./apiConfig";
import { ApiFetch } from "./apiFetch";

export interface TokenResponse {
  token: string;
  expiresAt: number;
}

export const authApi = {
  create: async (loginToken: string) => {
    const res = await ApiFetch.fetch<TokenResponse>({
      method: "POST",
      route: "/api/auth/token",
      headers: { Authorization: `Bearer ${loginToken}` },
    });
    ApiConfig.setToken(res.token);
    sessionStorage.setItem("token", res.token);

    return res;
  },
};
