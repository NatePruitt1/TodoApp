const BASE_URL = import.meta.env.VITE_API_URL;

export type ErrorCodeType = "AUTH" | "VAL" | "DB";

export interface ErrorCode {
    type: ErrorCodeType;
    code: string;
}

export class ApiError extends Error {
    constructor(status: number, error_code: string, issue: string) {
        super(issue);
        this.status = status;
        this.issue = issue;

        const errorCodeSplit = error_code.split("-")
        this.error_code = {
            type: errorCodeSplit[0] as ErrorCodeType,
            code: errorCodeSplit[1],
        }

    }
    status: number;
    issue: string;
    error_code: ErrorCode;
}

interface ApiEnvelope<T> {
    status: string;
    data: T;
}

interface ApiResponse<T> {
    data: T;
    headers: Headers;
}

export async function apiRequest<T>(path: string, jwt: string | null, init: RequestInit = {}): Promise<ApiResponse<T>> {
    const headers: HeadersInit = {
        "Content-Type": "application/json",
        ...(jwt ? { Authorization: "Bearer " + jwt} : {}),
        ...init.headers,
    }

    const resp = await fetch(BASE_URL + path, {...init, headers, credentials: "include"});
    if (!resp.ok) {
        const body = await resp.json().catch(() => ({}));
        const apiErr = new ApiError(resp.status, body.code, body.issue ?? "request_failed");
        console.error(apiErr.error_code)
        throw apiErr
    }

    if (resp.status === 204) {
        return {data: undefined as T, headers: resp.headers};
    } else {
        const body: ApiEnvelope<T> = await resp.json();
        return {data: body.data, headers: resp.headers}
    }
} 

export async function apiFetch<T>(path: string, jwt: string | null, init: RequestInit = {}): Promise<T>{
    const {data} = await apiRequest<T>(path, jwt, init)
    return data;
}