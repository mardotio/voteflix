import { format } from "date-fns";

import type { MovieStatus } from "../../utils/client/moviesApi";
import { MOVIE_STATUS_LABELS } from "../../utils/statusLabels";
import { Avatar } from "../Avatar";
import styles from "./MovieCard.module.scss";

interface MovieCardProps {
  movie: {
    id: string;
    name: string;
    status: MovieStatus;
    creator: {
      name: string;
      avatarUrl: string | null;
    };
    createdAt: number;
    watchedAt: null | number;
  };
  statusFilter: MovieStatus | "all" | null;
  onClick: () => void;
}

export const MovieCard = ({ movie, statusFilter, onClick }: MovieCardProps) => {
  return (
    <li className={styles.movie}>
      <button onClick={onClick}>
        <div className={styles.header}>
          <div className={styles.creator}>
            <Avatar
              name={movie.creator.name}
              src={movie.creator.avatarUrl}
              size={24}
            />
            <p>{movie.creator.name}</p>
          </div>
          <p className={styles.status}>{MOVIE_STATUS_LABELS[movie.status]}</p>
        </div>
        <h6 className={styles.name}>{movie.name}</h6>
        <p className={styles.date}>
          {format(
            new Date(
              statusFilter === "watched"
                ? (movie.watchedAt ?? movie.createdAt)
                : movie.createdAt,
            ),
            "P",
          )}
        </p>
      </button>
    </li>
  );
};
