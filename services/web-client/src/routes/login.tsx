import { useQuery } from "@tanstack/react-query";
import { Navigate, createFileRoute, useSearch } from "@tanstack/react-router";

import { useCurrentUserContext } from "../hooks/useCurrentUser";
import { ApiConfig, authApi, usersApi } from "../utils/client";

interface LoginLayoutSearchParams {
  t?: string;
}

const LoginLayout = () => {
  const { setCurrentUser } = useCurrentUserContext();
  const { t: token } = useSearch({ from: "/login" });
  const { data } = useQuery({
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

  if (data) {
    return (
      <Navigate to="/$serverId" params={{ serverId: data.list.id }} replace />
    );
  }

  return (
    <div
      style={{
        display: "flex",
        width: "100%",
        height: "100%",
        justifyContent: "center",
        alignItems: "center",
        flexDirection: "column",
      }}
    >
      <img
        style={{ width: "180px", height: "180px" }}
        src="/voteflix.svg"
        alt="Voteflix logo - three circles spread out from top left corner to bottom right corner"
      />
    </div>
  );
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
