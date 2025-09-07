class Config {
  baseEndpoint: string | null = null;
  headers: RequestInit["headers"];

  private _token: string | null = null;

  init() {
    const storedToken = sessionStorage.getItem("token");

    if (storedToken) {
      this.setToken(storedToken);
    }
  }

  setToken(t: string) {
    this._token = t;
    sessionStorage.setItem("token", t);
  }

  hasToken() {
    return this._token !== "" && this._token !== null;
  }

  getHeaders(): HeadersInit {
    return {
      ...(this.hasToken()
        ? {
            Authorization: `Bearer ${this._token}`,
          }
        : {}),
      ...(this.headers || {}),
    };
  }
}

export const ApiConfig = new Config();
