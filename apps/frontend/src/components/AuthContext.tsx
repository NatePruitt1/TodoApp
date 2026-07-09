import { createContext, useContext, useState, type ReactNode } from "react";

interface UserContextType {
    username: string | null;
    setUsername: (name: string | null) => void
}

const AuthContext = createContext<UserContextType | null>(null);

interface ProviderProps {
    children: ReactNode;
}

export function AuthProvider({ children }: ProviderProps) {
    const [username, setUsername] = useState<string | null>(null);
    const setUsernameWrap = (name: string | null) => {
        console.log(name);
        setUsername(name)
    }
    return (
        <AuthContext.Provider value= {{ username, setUsername: setUsernameWrap }}>
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