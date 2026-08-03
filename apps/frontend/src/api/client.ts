const BASE_URL = import.meta.env.VITE_API_URL;

export class ApiError extends Error {
    constructor(status: number, issue: string) {
        super(issue);
        this.status = status;
        this.issue = issue;
    }
    status: number;
    issue: string;
}

interface ApiEnvelope<T> {
    status: string;
    data: T;
}

export async function apiFetch<T>(path: string, jwt: string | null, init: RequestInit = {}): Promise<T>{
    const headers: HeadersInit = {
        "Content-Type": "application/json",
        ...(jwt ? { Authorization: "Bearer " + jwt} : {}),
        ...init.headers,
    }

    const resp = await fetch(BASE_URL + path, {...init, headers, credentials: "include"});
    if (!resp.ok) {
        const body = await resp.json().catch(() => ({}));
        throw new ApiError(resp.status, body.issue ?? "request_failed");
    }

    if (resp.status === 204) {
        return undefined as T;
    } else {
        const body: ApiEnvelope<T> = await resp.json();
        return body.data
    }
}