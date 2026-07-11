import { createContext, useContext, useEffect, useState, type PropsWithChildren } from "react";

export interface UserContextType {
    username: string;
    jwt: string;
    expiry: Date;
    setExpiry: (expiry: Date) => void
    setUsername: (name: string) => void
    setJWT: (jwt: string) => void
}

const AuthContext = createContext<UserContextType | null>(null);

export function AuthProvider({ children }: PropsWithChildren) {
    const [username, setUsername] = useState<string>("");
    const [jwt, setJWT] = useState<string>("");
    const [expiry, setExpiry] = useState<Date>(new Date());

    useEffect(() => {
        fetch("/refresh")
    })
    
    return (
        <AuthContext.Provider value= {{ username, jwt, expiry, setUsername: setUsername, setJWT: setJWT, setExpiry: setExpiry}}>
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