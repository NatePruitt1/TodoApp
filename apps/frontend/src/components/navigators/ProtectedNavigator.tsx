import { Outlet, Navigate } from "react-router";
import { useAuth } from "../AuthContext";
import { requireAuth } from "../../utils/Auth";

const BASE_URL = import.meta.env.VITE_API_URL;

export function ProtectedNavigator() {
    const userContext = useAuth();

    if (userContext.isLoading) {
        return null;
    }

    try {
        const s = requireAuth(userContext);
        console.log(s)
    } catch (e) {
        console.log(e)
        return (<><Navigate to="/login" replace={true} /></>)
    }

    return (
        <>
            <Outlet />
        </>
    )
}