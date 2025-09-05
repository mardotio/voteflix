import { ApiConfig } from "./apiConfig";

export interface ApiErrorData {
  status: number;
  statusText: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  body: any;
}

export class ApiFetch {
  static async generateError(
    response: Response,
    asText = false,
  ): Promise<ApiErrorData> {
    return {
      status: response.status,
      statusText: response.statusText,
      body: asText ? await response.text() : await response.json(),
    };
  }

  static getEndpoint(route: string) {
    if (ApiConfig.baseEndpoint) {
      return `${ApiConfig.baseEndpoint}${route}`;
    }
    return route;
  }

  static async fetch<Res, Body extends object | undefined = undefined>({
    route,
    method,
    headers = {},
    body,
  }: {
    route: string;
    method: Request["method"];
    body?: Body;
    headers?: Record<string, string>;
  }) {
    const response = await fetch(this.getEndpoint(route), {
      method,
      headers: {
        ...(ApiConfig.headers || {}),
        "content-type": "application/json",
        ...headers,
      },
      body: body ? JSON.stringify(body) : undefined,
    });

    console.log({ response });

    if (!response.ok) {
      const errorResponse = await this.generateError(
        response,
        !response.headers.get("content-type")?.includes("application/json"),
      );
      return errorResponse;
    }

    if (response.status === 204) {
      return undefined as unknown as Res;
    }

    const bodyResponse: Res = await response.json();
    return bodyResponse;
  }
}
