import { Outlet, useNavigate } from "react-router";
import { useAuth } from "../AuthContext";
import { requireAuth } from "../../utils/Auth";

export function ProtectedNavigator() {
    const navigate = useNavigate();
    const userContext = useAuth();

    try {
        requireAuth(userContext);
    } catch (e) {
        navigate('/login')
    }

    return (
        <>
            <Outlet />
        </>
    )
}