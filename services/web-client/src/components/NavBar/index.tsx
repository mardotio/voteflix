import { Link } from "@tanstack/react-router";

import { useCurrentUser } from "../../hooks/useCurrentUser";
import { MovieIcon } from "../Icon";
import { MovieIconFilled } from "../Icon/MovieIconFilled";
import styles from "./NavBar.module.scss";

export const NavBar = () => {
  const currentUser = useCurrentUser();
  return (
    <nav className={styles.main}>
      <ul className={styles.navigation}>
        <li className={styles.selected}>
          <Link to="/$serverId" params={{ serverId: currentUser.list.id }}>
            {(p) => {
              if (p.isActive) {
                return <MovieIconFilled />;
              }
              return <MovieIcon />;
            }}
          </Link>
        </li>
        <li>Home</li>
        <li>Home</li>
      </ul>
    </nav>
  );
};
