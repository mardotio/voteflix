import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { format } from "date-fns";

import { Avatar } from "../../../components/Avatar";
import { moviesApi } from "../../../utils/client/moviesApi";
import styles from "./movies.module.scss";

const MoviesLayout = () => {
  const { data } = useQuery({
    queryFn: async ({ queryKey }) => {
      return await moviesApi.listMovies(queryKey[1]);
    },
    queryKey: ["movies", { limit: 10, direction: "desc" }] as const,
  });

  return (
    <ul>
      {data?.data.map((m) => (
        <li key={m.id} className={styles.movie}>
          <div>
            <div className={styles.header}>
              <div className={styles.creator}>
                <Avatar
                  name={m.creator.name}
                  src={m.creator.avatarUrl}
                  size={24}
                />
                <p>{m.creator.name}</p>
              </div>
              <p className={styles.status}>{m.status}</p>
            </div>
            <h6 className={styles.name}>{m.name}</h6>
            <p className={styles.date}>{format(new Date(m.createdAt), "P")}</p>
          </div>
        </li>
      ))}
    </ul>
  );
};

export const Route = createFileRoute("/_loggedInLayout/$serverId/movies")({
  component: MoviesLayout,
});
