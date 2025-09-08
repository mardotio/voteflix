import { useQuery } from "@tanstack/react-query";
import { Navigate, createFileRoute, useSearch } from "@tanstack/react-router";

import { useCurrentUserContext } from "../hooks/useCurrentUser";
import { ApiConfig, authApi, isApiError, usersApi } from "../utils/client";

interface LoginLayoutSearchParams {
  t?: string;
}

const LoginLayout = () => {
  const { setCurrentUser } = useCurrentUserContext();
  const { t: token } = useSearch({ from: "/login" });
  const { data, status, error } = useQuery({
    queryFn: async () => {
      if (!token && !ApiConfig.hasToken()) {
        throw new Error("No token provided");
      }

      if (token) {
        await authApi.create(token);
      }

      const res = await usersApi.whoAmI();
      setCurrentUser(res);
      return res;
    },
    queryKey: ["login", { token }],
    retry: false,
  });

  if (status === "pending") {
    return <div>Logging in</div>;
  }

  if (data) {
    return (
      <Navigate to="/$serverId" params={{ serverId: data.list.id }} replace />
    );
  }

  if (!token) {
    return <div>Please provide token to log in</div>;
  }

  if (isApiError(error)) {
    return (
      <div>
        <pre>{JSON.stringify(error.body, null, 2)}</pre>
      </div>
    );
  }

  return <div>{error?.message}</div>;
};

export const Route = createFileRoute("/login")({
  component: LoginLayout,
  beforeLoad: () => {
    ApiConfig.init();
  },
  validateSearch: (
    search: Record<string, unknown>,
  ): LoginLayoutSearchParams => ({
    t:
      search.t && typeof search.t === "string" && search.t !== ""
        ? search.t
        : undefined,
  }),
});
