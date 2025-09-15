import { ApiFetch } from "./apiFetch";

export interface WhoAmIResponse {
  id: string;
  displayName: string;
  avatarUrl: string | null;
  list: {
    id: string;
    name: string;
    serverId: string;
    avatarId: string | null;
  };
}

export const usersApi = {
  whoAmI: async () =>
    ApiFetch.fetch<WhoAmIResponse>({
      route: "/api/users/whoami",
    }),
};
