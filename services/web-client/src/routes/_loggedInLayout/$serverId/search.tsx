import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";

import { MovieCard } from "../../../components/MovieCard";
import { MovieDetails } from "../../../components/MovieDetails";
import { moviesApi } from "../../../utils/client/moviesApi";
import styles from "./search.module.scss";

const SearchLayout = () => {
  const [query, setQuery] = useState("");
  const [debouncedQuery, setDebouncedQuery] = useState("");
  const [selectedMovie, setSelectedMovie] = useState<string | null>(null);

  const { data } = useQuery({
    queryFn: ({ queryKey }) =>
      moviesApi.listMovies({
        query: queryKey[1].query,
        limit: 10,
        direction: "desc",
      }),
    queryKey: ["movies", { query: debouncedQuery.trim() }] as const,
    enabled: debouncedQuery.trim().length > 0,
  });

  const { data: movieDetails } = useQuery({
    queryFn: ({ queryKey }) => moviesApi.getMovie(queryKey[1].id),
    queryKey: ["movie", { id: selectedMovie as string }] as const,
    enabled: selectedMovie !== null,
  });

  useEffect(() => {
    const t = setTimeout(() => setDebouncedQuery(query), 350);

    return () => clearTimeout(t);
  }, [query]);

  return (
    <>
      <div>
        <div className={styles.search}>
          <input
            placeholder="Find a movie..."
            name="movie"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
        <ul>
          {data?.data.map((m) => (
            <MovieCard
              key={m.id}
              movie={m}
              onClick={() => setSelectedMovie(m.id)}
              statusFilter={null}
            />
          ))}
        </ul>
      </div>
      <MovieDetails
        movie={movieDetails ?? null}
        onClose={() => setSelectedMovie(null)}
      />
    </>
  );
};

export const Route = createFileRoute("/_loggedInLayout/$serverId/search")({
  component: SearchLayout,
});
