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