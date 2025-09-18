import { ApiFetch } from "./apiFetch";

export interface CreateListResponse {
  id: string;
}

export interface CreateListRequest {
  discordUserId: string;
  discordUsername: string;
  discordServerId: string;
  discordServerName: string;
  discordServerAvatarId: string | null;
  discordNickname: string | null;
  discordAvatarId: string | null;
}

export interface PickMovieResponse {
  id: string;
  name: string;
  creator: {
    name: string;
    avatarUrl: string | null;
  };
  createdAt: string;
}

// eslint-disable-next-line import/prefer-default-export
export const BotApi = {
  createList: (payload: CreateListRequest) =>
    ApiFetch.fetch<CreateListResponse, CreateListRequest>(
      "/bot/list",
      "POST",
      payload,
    ),
  pickMovie: (serverId: string) =>
    ApiFetch.fetch<PickMovieResponse>(`/bot/${serverId}/movies/pick`, "GET"),
};
