import { useNavigate } from "react-router";
import { requireAuth } from "../utils/Auth";
import { useAuth } from "./AuthContext";
import { Header } from "./Header";

export function ProfileScreen() {
    const authContent = useAuth()
    const navigate = useNavigate();

    const logout = () => {
        try {
            const status = requireAuth(authContent)
            if(status.ok) {
                authContent.setJWT("")
                authContent.setUsername("")
                authContent.setExpiry(new Date(Date.now()))
            }
        } catch(error) {
            console.error(error)
        } finally {
            navigate("/login")
        }
    }

    return (
        <div>
            <Header title="Profile" />
            <button onClick={logout}>Log out</button>
        </div>
    )
}