import { Link } from "@tanstack/react-router";
import { useState } from "react";

import { useCurrentUser } from "../../hooks/useCurrentUser";
import { AddMovie } from "../AddMovie";
import { HomeIcon, MagnifyingGlassIcon, PlusSquareIcon } from "../Icon";
import styles from "./NavBar.module.scss";

export const NavBar = () => {
  const currentUser = useCurrentUser();
  const [isAddOpen, setIsAddOpen] = useState(false);

  return (
    <>
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
            <button onClick={() => setIsAddOpen(true)}>
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
      <AddMovie isOpen={isAddOpen} onClose={() => setIsAddOpen(false)} />
    </>
  );
};
