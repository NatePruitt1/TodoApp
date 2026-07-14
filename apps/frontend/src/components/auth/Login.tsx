import React, { useState } from 'react'
import './Login.css'
import { useNavigate } from 'react-router';
import { useAuth } from '../AuthContext';
import { jwtDecode } from 'jwt-decode';
import { type JwtData } from '../../types/Jwt';

const BASE_URL = import.meta.env.VITE_API_URL;

function Login() {
    const [formData, setFormData] = useState({username: "", password: ""});
    const userContext = useAuth();

    const navigate = useNavigate();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

    const handleSubmit = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();

        try {
            console.log(BASE_URL)
            let resp = await fetch(BASE_URL + "/api/v0/login",
                {
                    method: "POST",
                    body: JSON.stringify(formData),
                    credentials: "include"
                }
            )

            if(resp.ok) {
                const authHeader = resp.headers.get("Authorization")
                const token = authHeader && authHeader.replace(/^Bearer\s+/i, '');
                
                if(!token) {
                    throw new Error("Token not included in response.")
                }

                userContext.setUsername(formData.username);
                userContext.setJWT(token)

                const payload = jwtDecode<JwtData>(token)

                userContext.setExpiry(new Date(payload.exp * 1000))

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
                <input id="username-input" 
                    value={formData.username} 
                    onChange={handleChange} 
                    type='text'
                    name='username'
                    placeholder="Enter username"></input>
                <input 
                    id="password-input" 
                    value={formData.password} 
                    onChange={handleChange} 
                    type='password'
                    name='password'
                    placeholder="Enter password"></input>
                <button type='submit'>Login</button>
            </form>
        </>
    )
}

export default Login