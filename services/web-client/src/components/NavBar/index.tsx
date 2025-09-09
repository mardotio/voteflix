import { Link } from "@tanstack/react-router";

import { useCurrentUser } from "../../hooks/useCurrentUser";
import { HomeIcon, MagnifyingGlassIcon, PlusSquareIcon } from "../Icon";
import styles from "./NavBar.module.scss";

export const NavBar = () => {
  const currentUser = useCurrentUser();
  return (
    <nav className={styles.main}>
      <ul className={styles.navigation}>
        <li>
          <Link
            to="/$serverId/movies"
            params={{ serverId: currentUser.list.id }}
          >
            {(p) => {
              return (
                <HomeIcon
                  size={24}
                  iconStyle={p.isActive ? "solid" : "outline"}
                />
              );
            }}
          </Link>
        </li>
        <li className={styles["main-button"]}>
          <button>
            <PlusSquareIcon />
          </button>
        </li>
        <li>
          <Link
            to="/$serverId/search"
            params={{ serverId: currentUser.list.id }}
          >
            {(p) => {
              return (
                <MagnifyingGlassIcon
                  size={24}
                  iconStyle={p.isActive ? "solid" : "outline"}
                />
              );
            }}
          </Link>
        </li>
      </ul>
    </nav>
  );
};
