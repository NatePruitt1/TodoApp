import { useAuth } from "../components/AuthContext";
import { requireAuth } from "../utils/Auth";
import { apiFetch } from "./client";

export function useApi() {
    const authContext = useAuth();
    const makeAuthRequest = <T>(path: string, init?: RequestInit) => {
        const status = requireAuth(authContext);
        return apiFetch<T>(path, status.jwt, init)
    }

    return {
        get: <T>(path: string) => makeAuthRequest<T>(path),
        post: <T>(path: string, body: unknown) => makeAuthRequest<T>(path, {method:"POST", body: JSON.stringify(body)}),
        patch: <T>(path: string, body: unknown) => makeAuthRequest<T>(path, {method:"PATCH", body: JSON.stringify(body)}),
        delete: <T>(path: string) => makeAuthRequest<T>(path, {method:"DELETE"}),
    }
}