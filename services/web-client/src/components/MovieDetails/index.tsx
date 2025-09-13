import { useMutation, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import { useEffect, useState } from "react";

import { useCurrentUser } from "../../hooks/useCurrentUser";
import { type GetMovieResponse, moviesApi } from "../../utils/client/moviesApi";
import { Avatar } from "../Avatar";
import { Drawer } from "../Drawer";
import { LikeIcon } from "../Icon";
import styles from "./MovieDetails.module.scss";
import { MovieDetailsTimeline } from "./MovieDetailsTimeline";

interface MovieDetailsContentProps {
  movie: GetMovieResponse;
  view: "details" | "vote";
  setView: (v: "details" | "vote") => void;
}

const MovieDetailsContent = ({
  movie,
  view,
  setView,
}: MovieDetailsContentProps) => {
  const queryClient = useQueryClient();
  const currentUser = useCurrentUser();
  const userVote =
    movie.votes.find((v) => v.userId === currentUser.id)?.isApproval ?? null;
  const [vote, setVote] = useState<boolean | null>(userVote);
  const [isAddingVote, setIsAddingVote] = useState(false);
  const { mutate: addMovieVote, isPending } = useMutation({
    mutationFn: ({ approve, movieId }: { movieId: string; approve: boolean }) =>
      moviesApi.addMovieVote(movieId, approve),
    mutationKey: ["votes"],
    onSuccess: (p) => {
      queryClient.invalidateQueries({ queryKey: ["movie", { id: p.movieId }] });
      setView("details");
    },
  });

  const currentYear = new Date().getFullYear();

  if (view === "vote") {
    return (
      <div className={styles["vote-view"]}>
        <h6>Are you interested?</h6>
        <div className={styles["vote-group"]}>
          <button
            className={vote === true ? styles.active : undefined}
            onClick={() => setVote(true)}
          >
            <LikeIcon size={32} />
          </button>
          <button
            className={vote === false ? styles.active : undefined}
            onClick={() => setVote(false)}
          >
            <LikeIcon size={32} className={styles.invert} />
          </button>
        </div>

        <button
          disabled={vote === null || isPending}
          onClick={() => {
            if (vote === null || userVote === vote) {
              setView("details");
              return;
            }
            addMovieVote({ approve: vote, movieId: movie.id });
          }}
        >
          Submit
        </button>
      </div>
    );
  }

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
          watchedAt={movie.watchedAt}
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

export interface MovieDetailsProps {
  movie: GetMovieResponse | null;
  onClose: () => void;
}

export const MovieDetails = ({ movie, onClose }: MovieDetailsProps) => {
  const currentUser = useCurrentUser();
  const [view, setView] = useState<"details" | "vote">("details");
  const userVote =
    movie?.votes.find((v) => v.userId === currentUser.id)?.isApproval ?? null;

  useEffect(() => {
    if (!movie?.status) {
      return;
    }

    if (movie.status === "pending" && userVote === null) {
      setView("vote");
    } else {
      setView("details");
    }
  }, [movie?.status, userVote]);

  return (
    <Drawer
      height={view === "vote" ? "250px" : undefined}
      isOpen={movie !== null}
      onClose={onClose}
      className={styles["movie-details"]}
      header={movie?.name}
    >
      {movie && (
        <MovieDetailsContent movie={movie} view={view} setView={setView} />
      )}
    </Drawer>
  );
};
