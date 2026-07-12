import { createContext, useContext, useEffect, useRef, useState, type PropsWithChildren } from "react";
import type { UserRequestDTO } from "../types/UserDTO";
import { jwtDecode } from "jwt-decode";
import { type JwtData } from "../types/Jwt";

const BASE_URL = import.meta.env.VITE_API_URL;

export interface UserContextType {
    username: string;
    jwt: string;
    expiry: Date;
    isLoading: boolean;
    setLoading: (loading: boolean) => void
    setExpiry: (expiry: Date) => void
    setUsername: (name: string) => void
    setJWT: (jwt: string) => void
}

const AuthContext = createContext<UserContextType | null>(null);

export function AuthProvider({ children }: PropsWithChildren) {
    const [username, setUsername] = useState<string>("");
    const [jwt, setJWT] = useState<string>("");
    const [expiry, setExpiry] = useState<Date>(new Date());
    const [loading, setLoading] = useState<boolean>(true);

    const refreshSent = useRef(false);

    useEffect(() => {
        if(!refreshSent.current) {
            refreshSent.current = true;
            setLoading(true)
            fetch(BASE_URL + "/api/v0/refresh", {
                method: "POST",
                credentials: "include",
            }).then(async (resp) => {
                if(resp.ok) {
                    const authHeader = resp.headers.get("Authorization")
                    const token = authHeader && authHeader.replace(/^Bearer\s+/i, '');
                    
                    if(!token) {
                        throw new Error("Token not included in response.")
                    }

                    let bodyData: UserRequestDTO = await resp.json()
                    setUsername(bodyData.data.username)
                    setJWT(token)
                    
                    const jwtData = jwtDecode<JwtData>(token)
                    setExpiry(new Date(jwtData.exp * 1000))
                }
            }).finally(() => setLoading(false))
        }
    }, []);
    
    return (
        <AuthContext.Provider value= {{ username, jwt, expiry, isLoading:loading, setLoading: setLoading, setUsername: setUsername, setJWT: setJWT, setExpiry: setExpiry}}>
            {children}
        </AuthContext.Provider>
    )
}

export function useAuth() {
    const context = useContext(AuthContext);

    if(!context) {
        throw new Error("useAuth must be used within a AuthProvider element.");
    }

    return context;
}