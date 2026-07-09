import React, { useState } from 'react'
import './Login.css'
import { useNavigate } from 'react-router';

const BASE_URL = import.meta.env.VITE_API_URL;

function Login() {
    const [formData, setFormData] = useState({username: "", password: ""});

    const navigate = useNavigate();

    const handleSubmit = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();

        try {
            console.log(BASE_URL)
            let resp = await fetch(BASE_URL + "/api/v0/login",
                {
                    method: "POST",
                    body: JSON.stringify(formData)
                }
            )

            if(resp.ok) {
                setFormData({username: "", password: ""})
                navigate("/")
            } else {
                console.log(await resp.json())
                setFormData({username: "", password: ""})
                alert("Login failed.")
            }
        } catch(e) {
            console.log(e)
            setFormData({username: "", password: ""})
            alert("Login failed.")
        }
    };

    return (
        <>
            <form id="login-form" className="auth-form" onSubmit={handleSubmit}>
                <h1>Log in to Todo-App</h1>
                <input id="username-input" placeholder="Enter username"></input>
                <input id="password-input" placeholder="Enter password"></input>
                <button type='submit'>Login</button>
            </form>
        </>
    )
}

export default Login