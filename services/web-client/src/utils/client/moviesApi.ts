import { ApiFetch } from "./apiFetch";

export type MovieStatus = "watched" | "approved" | "rejected" | "pending";

export interface ListMoviesRequest {
  direction: "asc" | "desc";
  limit: number;
  status?: MovieStatus;
  before?: string;
  after?: string;
  query?: string;
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

interface MovieDetailsRating {
  rating: number;
  userId: string;
  createdAt: number;
  updatedAt: number | null;
}

export interface MovieDetailsVote {
  isApproval: boolean;
  userId: string;
  createdAt: number;
  updatedAt: number | null;
}

export interface MovieDetailsUsersMap {
  [UserId: string]: {
    name: string;
    avatarUrl: string | null;
  };
}

export interface GetMovieResponse {
  id: string;
  name: string;
  listId: string;
  status: MovieStatus;
  votes: MovieDetailsVote[];
  ratings: MovieDetailsRating[];
  creatorId: string;
  createdAt: number;
  updatedAt: number | null;
  watchedAt: number | null;
  users: MovieDetailsUsersMap;
}

export interface AddMovieVoteRequest {
  approve: boolean;
}

export interface AddMovieRatingRequest {
  rating: number;
}

export interface AddMovieVoteResponse {
  movieId: string;
  approved: boolean;
  createdAt: number;
  updatedAt: number | null;
}

export interface UpdateMovieRequest {
  name?: string;
  status?: MovieStatus;
}

export interface UpdateMovieResponse {
  id: string;
  name: string;
  status: MovieStatus;
  createdAt: number;
  updatedAt: number | null;
  watchedAt: number | null;
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
  getMovie: (movieId: string) =>
    ApiFetch.fetch<GetMovieResponse>({
      method: "GET",
      route: `/api/movies/${movieId}`,
    }),
  addMovieVote: (movieId: string, approve: boolean) =>
    ApiFetch.fetch<AddMovieVoteResponse, AddMovieVoteRequest>({
      method: "PUT",
      route: `/api/movies/${movieId}/votes`,
      body: { approve },
    }),
  addMovieRating: (movieId: string, rating: number) =>
    ApiFetch.fetch<AddMovieVoteResponse, AddMovieRatingRequest>({
      method: "PUT",
      route: `/api/movies/${movieId}/ratings`,
      body: { rating },
    }),
  updateMovie: (movieId: string, body: UpdateMovieRequest) =>
    ApiFetch.fetch<UpdateMovieResponse, UpdateMovieRequest>({
      method: "PATCH",
      route: `/api/movies/${movieId}`,
      body,
    }),
};
