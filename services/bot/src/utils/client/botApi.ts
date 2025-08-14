import {ApiFetch} from './apiFetch';

export interface CreateListResponse {
    id: string;
}

export interface CreateListRequest {
    discordUserId: string;
    discordUsername: string;
    discordServerId: string;
    discordServerName: string;
    discordNickname: string | null;
    discordAvatarId: string | null;
}

// eslint-disable-next-line import/prefer-default-export
export const BotApi = {
    createList: (payload: CreateListRequest) => ApiFetch.fetch<CreateListResponse, CreateListRequest>('/bot/list', 'POST', payload),
};