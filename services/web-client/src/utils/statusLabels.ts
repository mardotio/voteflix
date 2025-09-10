import type { MovieStatus } from "./client/moviesApi";

export const MOVIE_STATUS_LABELS: Record<MovieStatus, string> = {
  pending: "Suggested",
  approved: "Approved",
  watched: "Watched",
  rejected: "Rejected",
};
