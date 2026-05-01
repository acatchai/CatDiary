const TOKEN_KEY = 'catdiary.token'

export function getToken() {
    return localStorage.getItem(TOKEN_KEY) || ''
}

export function setToken(token) {
    if (typeof token === 'string' && token.trim() !== '') {
        localStorage.setItem(TOKEN_KEY, token)
        return
    }

    localStorage.removeItem(TOKEN_KEY)
}
