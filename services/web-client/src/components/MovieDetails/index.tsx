import { format } from "date-fns";

import type { GetMovieResponse } from "../../utils/client/moviesApi";
import { Avatar } from "../Avatar";
import styles from "./MovieDetails.module.scss";
import { MovieDetailsTimeline } from "./MovieDetailsTimeline";

interface MovieDetailsProps {
  movie: GetMovieResponse | null;
}

export const MovieDetails = ({ movie }: MovieDetailsProps) => {
  if (!movie) {
    return null;
  }

  const currentYear = new Date().getFullYear();

  return (
    <>
      <div>
        <h5 className={styles["section-header"]}>Timeline</h5>
        <MovieDetailsTimeline
          createdAt={movie.createdAt}
          status={movie.status}
          creatorId={movie.creatorId}
          users={movie.users}
          votes={movie.votes}
          watchedAt={movie.updatedAt /*TODO: need watched at field from BE*/}
        />
      </div>
      {movie.ratings.length > 0 && (
        <div>
          <h5 className={styles["section-header"]}>Ratings</h5>
          <ul>
            {movie.ratings.map((r) => (
              <li key={r.userId} className={styles.rating}>
                <div className={styles["rating-header"]}>
                  <div className={styles["rating-header--user"]}>
                    <Avatar
                      size={24}
                      name={movie.users[r.userId].name}
                      src={movie.users[r.userId].avatarUrl}
                    />
                    <p>{movie.users[r.userId].name}</p>
                  </div>
                  <p className={styles["rating-header--timestamp"]}>
                    {format(
                      new Date(r.createdAt),
                      new Date(r.createdAt).getFullYear() === currentYear
                        ? "MMM d"
                        : "MMM d yyyy",
                    )}
                  </p>
                </div>
                <p>{r.rating} / 10</p>
              </li>
            ))}
          </ul>
        </div>
      )}
    </>
  );
};
