import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";

import { moviesApi } from "../../utils/client/moviesApi";
import { CloseIcon, PlusSquareIcon } from "../Icon";
import styles from "./AddMovie.module.scss";

interface AddMovieProps {
  isOpen: boolean;
  onClose: () => void;
}

export const AddMovie = ({ isOpen, onClose }: AddMovieProps) => {
  const queryClient = useQueryClient();
  const [movie, setMovie] = useState("");
  const [hasNew, setHasNew] = useState(false);
  const { mutate: createMovie, isPending } = useMutation({
    mutationFn: (name: string) => moviesApi.createMovie({ name }),
    mutationKey: ["movies"],
    onSuccess: () => {
      setMovie("");
      setHasNew(true);
    },
  });

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
    <div
      className={`${styles.container} ${isOpen ? styles.open : styles.closed}`}
    >
      <div className={styles.main}>
        <button onClick={close} className={styles["close-button"]}>
          <CloseIcon size={32} />
        </button>
        <div className={styles["movie-input"]}>
          <input
            placeholder="Jurassic Park"
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
      </div>
    </div>
  );
};
