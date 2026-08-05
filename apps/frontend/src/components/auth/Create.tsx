import React, { useState, useRef } from 'react'
import './Login.css'
import { Link, useNavigate } from 'react-router';
import { useAuth } from '../AuthContext';
import { jwtDecode } from 'jwt-decode';
import type { JwtData } from '../../types/Jwt';
import { type InputHandle, Input } from '../Input';
import { handleLoginFormError, validatePassword, validateUsername } from './validate';
import { usePublicApi } from '../../api/useApi';
import { ApiError } from '../../api/client';

function Create() {
    const [formData, setFormData] = useState({username: "", password: ""});
        const userContext = useAuth();
        const usernameInputRef = useRef<InputHandle>(null);
        const passwordInputRef = useRef<InputHandle>(null);
        const publicApi = usePublicApi();
        const navigate = useNavigate();
    
        const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
        };

        const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
            handleChange(e)
            const valid = validateUsername(e.target.value)
            if(!valid.ok) {
                usernameInputRef.current?.setError(valid.reason)
            }
        }

        const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
            handleChange(e)
            const valid = validatePassword(e.target.value)
            if(!valid.ok) {
                passwordInputRef.current?.setError(valid.reason)
            }
        }
    
        const handleSubmit = async (event: React.SubmitEvent<HTMLFormElement>) => {
            event.preventDefault();

            const validUsername = validateUsername(formData.username)
            const validPassword = validatePassword(formData.password)

            if(!validUsername.ok || !validPassword.ok) return;
        
            try {
                let resp = await publicApi.post('/auth/register', formData);

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
            } catch(e) {
                console.log(e)
                if(e instanceof ApiError && usernameInputRef.current && passwordInputRef.current) {
                    usernameInputRef.current.setError("")
                    passwordInputRef.current.setError("")
                    handleLoginFormError(e.error_code, usernameInputRef.current, passwordInputRef.current)
                }
            }
        };
    
        return (
            <>
                <form id="create-form" className="auth-form" onSubmit={handleSubmit}>
                    <h1>Create Account</h1>
                    <Input ref={usernameInputRef}
                        id="username-input" 
                        value={formData.username} 
                        onChange={handleUsernameChange} 
                        type='text'
                        name='username'
                        placeholder="Enter username" />
                    <Input ref={passwordInputRef} 
                        id="password-input" 
                        value={formData.password} 
                        onChange={handlePasswordChange} 
                        type='password'
                        name='password'
                        placeholder="Enter password" />
                    <button className="small-button" type='submit'>Create &rarr;</button>
                </form>
                <Link className="auth-link" to={{pathname: "/login"}}>Log in to your account.</Link>
            </>
        )
}

export default Create