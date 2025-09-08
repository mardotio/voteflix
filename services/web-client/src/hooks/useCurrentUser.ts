import { createContext, useContext } from "react";

import type { WhoAmIResponse } from "../utils/client";

interface CurrentUserContextData {
  currentUser: WhoAmIResponse | null;
  setCurrentUser: (v: WhoAmIResponse | null) => void;
}

export const CurrentUserContext = createContext<CurrentUserContextData>({
  currentUser: null,
  setCurrentUser: () => {},
});

export const useCurrentUserContext = () => useContext(CurrentUserContext);

export const useCurrentUser = () => {
  const { currentUser } = useCurrentUserContext();

  if (!currentUser) {
    throw new Error(
      "You should only use `userCurrentUser` from a logged in context",
    );
  }

  return currentUser;
};
