class Config {
  baseEndpoint: string | null = null;

  headers: RequestInit["headers"];
}

export const ApiConfig = new Config();
