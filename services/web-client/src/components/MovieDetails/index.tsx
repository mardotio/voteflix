import { useMutation, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import { useEffect, useState } from "react";

import { useCurrentUser } from "../../hooks/useCurrentUser";
import {
  type GetMovieResponse,
  type MovieStatus,
  moviesApi,
} from "../../utils/client/moviesApi";
import { Avatar } from "../Avatar";
import { Drawer } from "../Drawer";
import { LikeIcon } from "../Icon";
import { RatingSlider } from "../RatingSlider";
import styles from "./MovieDetails.module.scss";
import { MovieDetailsTimeline } from "./MovieDetailsTimeline";

type View = "details" | "vote" | "reaction";

interface MovieDetailsContentProps {
  movie: GetMovieResponse;
  view: View;
  setView: (v: View) => void;
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
  const userRating =
    movie.ratings.find((v) => v.userId === currentUser.id)?.rating ?? null;
  const [rating, setRating] = useState<number | null>(userRating);
  const [vote, setVote] = useState<boolean | null>(userVote);
  const { mutate: addMovieVote, isPending: isVoting } = useMutation({
    mutationFn: ({ approve, movieId }: { movieId: string; approve: boolean }) =>
      moviesApi.addMovieVote(movieId, approve),
    mutationKey: ["votes"],
    onSuccess: (res) => {
      queryClient.invalidateQueries({
        queryKey: ["movie", { id: res.movieId }],
      });
      setView("details");
    },
  });
  const { mutate: addMovieRating, isPending: isRating } = useMutation({
    mutationFn: ({ rating, movieId }: { movieId: string; rating: number }) =>
      moviesApi.addMovieRating(movieId, rating),
    mutationKey: ["reactions"],
    onSuccess: (res) => {
      queryClient.invalidateQueries({
        queryKey: ["movie", { id: res.movieId }],
      });
      setView("details");
    },
  });
  const { mutate: updateMovie, isPending: isChangingStatus } = useMutation({
    mutationFn: ({
      status,
      movieId,
    }: {
      movieId: string;
      status: MovieStatus;
    }) => moviesApi.updateMovie(movieId, { status }),
    mutationKey: ["movie"],
    onSuccess: (res) => {
      queryClient.invalidateQueries({ queryKey: ["movie", { id: res.id }] });
      queryClient.invalidateQueries({ queryKey: ["movies"] });
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
          disabled={vote === null || isVoting}
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

  if (view === "reaction") {
    return (
      <div className={styles["vote-view"]}>
        <h6>What did you think?</h6>
        <RatingSlider onChange={setRating} value={rating} />
        <button
          disabled={rating === null || isRating}
          onClick={() => {
            if (rating === null || userRating === rating) {
              setView("details");
              return;
            }
            addMovieRating({ rating, movieId: movie.id });
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
        {movie.status === "approved" && (
          <button
            className={styles["status-toggle"]}
            onClick={() =>
              updateMovie({ movieId: movie.id, status: "watched" })
            }
            disabled={isChangingStatus}
          >
            Mark as watched
          </button>
        )}
        {movie.status === "watched" && movie.ratings.length <= 0 && (
          <button
            className={styles["status-toggle"]}
            onClick={() =>
              updateMovie({ movieId: movie.id, status: "approved" })
            }
            disabled={isChangingStatus}
          >
            Mark as unwatched
          </button>
        )}
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
  const [view, setView] = useState<View>("details");
  const userVote =
    movie?.votes.find((v) => v.userId === currentUser.id)?.isApproval ?? null;
  const userRating =
    movie?.ratings.find((v) => v.userId === currentUser.id)?.rating ?? null;

  useEffect(() => {
    if (!movie?.status) {
      return;
    }

    if (movie.status === "pending" && userVote === null) {
      setView("vote");
    } else if (movie.status === "watched" && userRating === null) {
      setView("reaction");
    } else {
      setView("details");
    }
  }, [movie?.status, userVote, userRating]);

  return (
    <Drawer
      height={
        view === "details" ? undefined : view === "reaction" ? "280px" : "250px"
      }
      isOpen={movie !== null}
      onClose={onClose}
      className={styles["movie-details"]}
      header={movie?.name}
      onBack={view !== "details" ? () => setView("details") : undefined}
    >
      {movie && (
        <MovieDetailsContent movie={movie} view={view} setView={setView} />
      )}
    </Drawer>
  );
};
