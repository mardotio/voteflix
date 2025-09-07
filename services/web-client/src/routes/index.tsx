import { createFileRoute, redirect } from "@tanstack/react-router";

import { ApiConfig } from "../utils/client";

export const Route = createFileRoute("/")({
  component: Index,
  beforeLoad: () => {
    if (!ApiConfig.hasToken()) {
      throw redirect({
        to: "/login",
      });
    }
  },
});

function Index() {
  return (
    <div className="p-2">
      <h3>Welcome Home!</h3>
    </div>
  );
}
