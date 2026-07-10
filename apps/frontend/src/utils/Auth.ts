import type { UserContextType } from "../components/AuthContext";

export type AuthStatus = {ok: true, username: string, jwt: string; expiry: Date}
    | { ok: false; reason: "missing-token" | "expired-token" | "invalid-token"}

export function getAuthStatus(userContext: UserContextType): AuthStatus {
    if(userContext.jwt === "") {
        return {ok: false, reason: "missing-token"};
    } else if(!userContext.expiry || userContext.expiry.getTime() <= Date.now()) {
        return {ok: false, reason: "expired-token"};
    } else {
        return {ok: true, username: userContext.username, jwt: userContext.jwt, expiry: userContext.expiry};
    }
}

export function requireAuth(userContext: UserContextType) {
    const status = getAuthStatus(userContext);
    if(!status.ok) {
        throw new Error("Unauthenticated: " + status.reason);
    }
    return status;
}