import { ApiFetch } from "./apiFetch";

export const usersApi = {
  whoAmI: async () =>
    ApiFetch.fetch<object>({
      method: "GET",
      route: "/api/users/whoami",
    }),
};
