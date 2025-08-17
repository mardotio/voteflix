import {ApiConfig} from './client';
import {BOT_ENVIRONMENT} from "./environment";

ApiConfig.baseEndpoint = 'http://api:8000';
ApiConfig.headers = {
    "Bot-Api-Key": BOT_ENVIRONMENT.API_BOT_KEY,
}

export * from './client';