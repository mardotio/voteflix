import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";

import { moviesApi } from "../../utils/client/moviesApi";
import { MOVIE_STATUS_LABELS } from "../../utils/statusLabels";
import { Drawer, type DrawerProps } from "../Drawer";
import { PlusSquareIcon } from "../Icon";
import styles from "./AddMovie.module.scss";
import { MOVIES } from "./movies";

export const AddMovie = ({
  isOpen,
  onClose,
}: Omit<DrawerProps, "children">) => {
  const queryClient = useQueryClient();
  const [movie, setMovie] = useState("");
  const [hasNew, setHasNew] = useState(false);
  const [debouncedQuery, setDebouncedQuery] = useState("");
  const [placeholderIndex, setPlaceHolderIndex] = useState(0);

  const { mutate: createMovie, isPending } = useMutation({
    mutationFn: (name: string) => moviesApi.createMovie({ name }),
    mutationKey: ["movies"],
    onSuccess: () => {
      setMovie("");
      setHasNew(true);
    },
  });
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

  useEffect(() => {
    const t = setTimeout(() => setDebouncedQuery(movie), 350);

    return () => clearTimeout(t);
  }, [movie]);

  useEffect(() => {
    if (!isOpen || movie.trim()) {
      return;
    }

    const t = setInterval(
      () => setPlaceHolderIndex((i) => (i + 1) % MOVIES.length),
      1500,
    );

    return () => {
      clearInterval(t);
    };
  }, [isOpen, movie]);

  const submit = () => {
    const sanitizedMovie = movie.trim();
    if (sanitizedMovie) {
      createMovie(sanitizedMovie);
    }
  };

  const close = () => {
    if (hasNew) {
      queryClient.invalidateQueries({ queryKey: ["movies"] });
      setHasNew(false);
    }
    setMovie("");
    onClose();
  };

  return (
    <Drawer onClose={close} isOpen={isOpen} className={styles.main}>
      {!!data && (
        <ul className={styles["related-movies"]}>
          {data.data.map((m) => (
            <li key={m.id}>
              <p className={styles.name}>{m.name}</p>
              <p className={styles.status}>{MOVIE_STATUS_LABELS[m.status]}</p>
            </li>
          ))}
        </ul>
      )}
      <div className={styles["movie-input"]}>
        <input
          placeholder={`I want to watch ${MOVIES[placeholderIndex]}`}
          name="movie"
          value={movie}
          onChange={(e) => setMovie(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              submit();
            }
          }}
          disabled={isPending}
        />
        <button onClick={submit}>
          <PlusSquareIcon size={32} />
        </button>
      </div>
    </Drawer>
  );
};
