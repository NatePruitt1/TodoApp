import { Outlet, Navigate } from "react-router";
import { useAuth } from "../AuthContext";
import { requireAuth } from "../../utils/Auth";

export function ProtectedNavigator() {
    const userContext = useAuth();

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