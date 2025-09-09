import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_loggedInLayout/$serverId/")({
  beforeLoad: () =>
    redirect({ from: "/$serverId", to: "./movies", replace: true }),
});
