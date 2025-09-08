import { ApiFetch } from "./apiFetch";

export interface WhoAmIResponse {
  id: string;
  displayName: string;
  avatarUrl: string | null;
  list: {
    id: string;
    name: string;
    serverId: string;
  };
}

export const usersApi = {
  whoAmI: async () =>
    ApiFetch.fetch<WhoAmIResponse>({
      method: "GET",
      route: "/api/users/whoami",
    }),
};
