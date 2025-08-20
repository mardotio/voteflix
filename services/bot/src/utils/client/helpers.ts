import { ApiError } from "./apiFetch";

export const isErrorResponse = <T, K = Exclude<T, ApiError>>(
  response: K | ApiError,
): response is ApiError => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const r = response as any;

  if (!r) {
    return false;
  }

  return r.status && r.statusText;
};

export const isSuccessResponse = <T, K = Exclude<T, ApiError>>(
  response: K | ApiError,
): response is K => !isErrorResponse(response);
