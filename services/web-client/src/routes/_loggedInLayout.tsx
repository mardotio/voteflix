import { Outlet, createFileRoute, redirect } from "@tanstack/react-router";

import { Avatar } from "../components/Avatar";
import { NavBar } from "../components/NavBar";
import { useCurrentUser } from "../hooks/useCurrentUser";
import { ApiConfig } from "../utils/client";
import styles from "./_loggedInLayout.module.scss";

export const Route = createFileRoute("/_loggedInLayout")({
  component: RouteComponent,
  beforeLoad: () => {
    if (!ApiConfig.hasToken()) {
      throw redirect({
        to: "/login",
      });
    }
  },
});

function RouteComponent() {
  const currentUser = useCurrentUser();

  return (
    <div className={styles.content}>
      <div className={styles.body}>
        <div className={styles["header-outer"]}>
          <div className={styles.header}>
            <Avatar
              name={currentUser.list.name}
              src={currentUser.list.avatarId}
            />
            <h3>{currentUser.list.name.toLocaleUpperCase()}</h3>
          </div>
        </div>

        <div className={styles["page-content"]}>
          <Outlet />
        </div>
      </div>

      <NavBar />
    </div>
  );
}
