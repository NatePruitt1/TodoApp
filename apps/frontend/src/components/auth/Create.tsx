import React, { useState } from 'react'
import './Login.css'
import { Link, useNavigate } from 'react-router';
import { useAuth } from '../AuthContext';
import { jwtDecode } from 'jwt-decode';
import type { JwtData } from '../../types/Jwt';

const BASE_URL = import.meta.env.VITE_API_URL;

function Create() {
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

            if(formData.password.trim().length < 6) {
                alert("Password too short.")
                return
            }

            if(formData.username.trim().length < 3) {
                alert("Username too short.")
                return
            }
        
            try {
                console.log(BASE_URL)
                let resp = await fetch(BASE_URL + "/auth/register",
                    {
                        method: "POST",
                        body: JSON.stringify(formData)
                    }
                )
    
                if(resp.ok) {
                    const authHeader = resp.headers.get("Authorization")
                    const token = authHeader && authHeader.replace(/^Bearer\s+/i, '');
                    
                    if(token == null) {
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
                    alert("Create failed.")
                }
            } catch(e) {
                console.log(e)
                setFormData({username: "", password: ""})
                alert("Create failed.")
            }
        };
    
        return (
            <>
                <form id="create-form" className="auth-form" onSubmit={handleSubmit}>
                    <h1>Create Account for Todo-App</h1>
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
                <Link to={{pathname: "/login"}}>Log in to your account.</Link>
            </>
        )
}

export default Create