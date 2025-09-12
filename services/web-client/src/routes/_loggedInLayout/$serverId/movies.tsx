import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";

import { Drawer } from "../../../components/Drawer";
import { CircleCaretIcon, SortIcon } from "../../../components/Icon";
import { MovieCard } from "../../../components/MovieCard";
import { MovieDetails } from "../../../components/MovieDetails";
import { type MovieStatus, moviesApi } from "../../../utils/client/moviesApi";
import { MOVIE_STATUS_LABELS } from "../../../utils/statusLabels";
import styles from "./movies.module.scss";

const LABELS: Record<MovieStatus | "all", string> = {
  all: "All",
  ...MOVIE_STATUS_LABELS,
};

const MoviesLayout = () => {
  const [status, setStatus] = useState<MovieStatus | "all">("all");
  const [direction, setDirection] = useState<"asc" | "desc">("desc");
  const [selectedMovie, setSelectedMovie] = useState<string | null>(null);

  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } =
    useInfiniteQuery({
      queryKey: [
        "movies",
        {
          limit: 10,
          direction: direction,
          status: status === "all" ? undefined : status,
        },
      ] as const,
      queryFn: async ({ queryKey, pageParam }) => {
        return await moviesApi.listMovies({
          ...queryKey[1],
          ...(pageParam ? { after: pageParam } : {}),
        });
      },
      getNextPageParam: (lastPage) => lastPage.next,
      initialPageParam: "",
      staleTime: 1000 * 60,
    });

  const { data: movieDetails } = useQuery({
    queryFn: ({ queryKey }) => moviesApi.getMovie(queryKey[1].id),
    queryKey: ["movie", { id: selectedMovie as string }] as const,
    enabled: selectedMovie !== null,
  });

  return (
    <>
      <div className={styles["list-options"]}>
        <label htmlFor="stuff" className={styles["status-selector"]}>
          <select
            id="stuff"
            onChange={(e) => setStatus(e.target.value as MovieStatus | "all")}
          >
            {Object.entries(LABELS).map(([k, v]) => (
              <option key={k} value={k}>
                {v}
              </option>
            ))}
          </select>
          <CircleCaretIcon size={20} />
        </label>
        <button
          className={direction === "desc" ? styles["dir-desc"] : undefined}
          onClick={() => setDirection((d) => (d === "asc" ? "desc" : "asc"))}
        >
          <SortIcon size={24} />
        </button>
      </div>
      <ul>
        {data?.pages.map((p) =>
          p.data.map((m) => (
            <MovieCard
              key={m.id}
              movie={m}
              onClick={() => setSelectedMovie(m.id)}
            />
          )),
        )}
      </ul>
      {hasNextPage && (
        <div className={styles["load-button"]}>
          <button
            onClick={() => fetchNextPage()}
            disabled={!hasNextPage || isFetchingNextPage}
          >
            {isFetchingNextPage ? "Loading more..." : "Load More"}
          </button>
        </div>
      )}
      <Drawer
        isOpen={selectedMovie !== null}
        onClose={() => setSelectedMovie(null)}
        className={styles["movie-details"]}
        header={movieDetails?.name}
      >
        <MovieDetails movie={movieDetails ?? null} />
      </Drawer>
    </>
  );
};

export const Route = createFileRoute("/_loggedInLayout/$serverId/movies")({
  component: MoviesLayout,
});
