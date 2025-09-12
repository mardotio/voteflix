import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";

import { moviesApi } from "../../utils/client/moviesApi";
import { Drawer, type DrawerProps } from "../Drawer";
import { PlusSquareIcon } from "../Icon";
import styles from "./AddMovie.module.scss";

export const AddMovie = ({
  isOpen,
  onClose,
}: Omit<DrawerProps, "children">) => {
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
    <Drawer onClose={close} isOpen={isOpen} className={styles.main}>
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
    </Drawer>
  );
};
