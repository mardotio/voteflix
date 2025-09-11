import { ApiFetch } from "./apiFetch";

export type MovieStatus = "watched" | "approved" | "rejected" | "pending";

export interface ListMoviesRequest {
  direction: "asc" | "desc";
  limit: number;
  status?: MovieStatus;
  before?: string;
  after?: string;
}

export interface Movie {
  id: string;
  name: string;
  status: MovieStatus;
  creator: {
    name: string;
    avatarUrl: string | null;
  };
  createdAt: number;
  updatedAt: number | null;
}

export interface ListMoviesResponse {
  data: Movie[];
  next: string | null;
  previous: string | null;
}

export interface CreateMovieResponse {
  id: string;
  listId: string;
  name: string;
  status: string;
  createdAt: number;
}

export interface CreateMovieRequest {
  name: string;
}

export const moviesApi = {
  listMovies: (options: ListMoviesRequest) =>
    ApiFetch.fetch<ListMoviesResponse>({
      method: "GET",
      route: "/api/movies",
      searchParams: options as unknown as Record<string, string>,
    }),
  createMovie: (body: CreateMovieRequest) =>
    ApiFetch.fetch<CreateMovieResponse, CreateMovieRequest>({
      method: "POST",
      route: "/api/movies",
      body,
    }),
};
