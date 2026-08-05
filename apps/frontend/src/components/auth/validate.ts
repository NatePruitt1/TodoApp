import type { ErrorCode } from "../../api/client"
import type { InputHandle } from "../Input"

export const validateUsername = (username: string):{ok: boolean, reason: string} => {
    const usernameRegex = /^[a-zA-Z0-9]+(?:[._-][a-zA-Z0-9]+)*$/
    
    if(!username || username.trim().length <= 3) {
        return {
            ok: false,
            reason: "Username must be longer than 3 characters."
        }
    }

    if(username.trim().length >= 16) {
        return {
            ok: false,
            reason: "Username must be shorter than 16 characters"
        }
    }

    if(!usernameRegex.test(username)) {
        return {
                ok: false,
                reason: "Username my only contain A-Z, 0-9, \".\", \"_\", and \"-\". You may not repeat special characters."
        }
    }

    return {
        ok: true,
        reason: ""
    }
}

export const validatePassword = (password: string): {ok: boolean, reason: string} => {
    const spaceRegex = /\s/
    const uppercaseRegex = /[A-Z]/
    const specialCharRegex = /[^A-Za-z0-9\s]/

    if(spaceRegex.test(password)) {
        return {
            ok: false,
            reason: "Password may not contain spaces."
        }
    }

    if(!password || password.trim().length < 6) {
        return {
            ok: false,
            reason: "Password must be longer than 6 characters."
        }
    }

    if(password.trim().length >= 256) {
        return {
            ok: false,
            reason: "Password must be shorter than 256 characters."
        }
    }

    if(!uppercaseRegex.test(password)) {
        return {
            ok: false,
            reason: "Password must have atleast 1 uppercase character."
        }
    }

    if(!specialCharRegex.test(password)) {
        return {
            ok: false,
            reason: "Password must have atleast 1 special character."
        }
    }

    return {
        ok: true,
        reason: ""
    }
}

export const handleLoginFormError = (error_code: ErrorCode, usernameInput: InputHandle, passwordHandle: InputHandle) => {
    if(error_code.type != "AUTH") {
        alert("Unexpect error type. Please try again.")
    } else {
        switch(error_code.code) {
            case "105":
                usernameInput.setError("Username or Password rejected. Please try again.")
                break;
            case "000":
                usernameInput.setError("Username may not contain spaces.")
                break;
            case "001":
                usernameInput.setError("Username my only contain A-Z, 0-9, \".\", \"_\", and \"-\". You may not repeat special characters.")
                break;
            case "002":
                usernameInput.setError("Username must be longer than 3 characters.")
                break;
            case "003":
                usernameInput.setError("Username must be shorter than 16 characters.")
                break;
            case "200":
                passwordHandle.setError("Password must not contain spaces.")
                break;
            case "201":
                passwordHandle.setError("Password must contain at least 1 special character.")
                break;
            case "202":
                passwordHandle.setError("Password must be longer than 6 characters.")
                break;
            case "203":
                passwordHandle.setError("Password must be shorter than 256 characters.")
                break;
            case "204":
                passwordHandle.setError("Password must contain an uppercase character.")
                break;
        }
    }
}