import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { useState } from "react";

import { CurrentUserContext } from "./hooks/useCurrentUser";
// Import the generated route tree
import { routeTree } from "./routeTree.gen";
import type { WhoAmIResponse } from "./utils/client";

// Create a new router instance
const router = createRouter({ routeTree });

// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const queryClient = new QueryClient();

export const App = () => {
  const [currentUser, setCurrentUser] = useState<WhoAmIResponse | null>(null);

  return (
    <CurrentUserContext value={{ currentUser, setCurrentUser }}>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </CurrentUserContext>
  );
};
