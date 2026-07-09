import { Outlet } from "react-router"
import "./AuthNavigator.css"

function AuthNavigator() {
    return (
        <>
            <div className="auth-layout">
                <div className="auth-card">
                    <Outlet />
                </div>
            </div>
        </>
    )
}

export default AuthNavigator