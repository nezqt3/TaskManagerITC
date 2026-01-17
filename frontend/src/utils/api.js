import { handleUnauthorized } from "./auth";

export async function apiFetch(url, options = {}) {
  const response = await fetch(url, options);
  if (response.status === 401) {
    handleUnauthorized();
  }
  return response;
}
