import { getToken } from './auth'

export async function apiRequest(path, options = {}) {
    const url = path.startsWith('http')
        ? path
        : `/api/v1${path.startsWith('/') ? path : `/${path}`}`

    const headers = new Headers(options.headers || {})
    if (!headers.has('Content-Type') && options.body != null) {
        headers.set('Content-Type', 'application/json')
    }

    const token = getToken()
    if (token) headers.set('Authorization', `Bearer ${token}`)

    const res = await fetch(url, {
        ...options,
        headers,
    })

    const rawText = await res.text().catch(() => '')
    const data = (() => {
        if (!rawText) return null
        try {
            return JSON.parse(rawText)
        } catch {
            return rawText
        }
    })()

    if (!res.ok) {
        const err = new Error(
            (typeof data === 'object' && data ? (data.error || data.message) : '') ||
            (typeof data === "string" ? data : '') ||
            `请求失败 (${res.status})`
        )
        err.status = res.status
        err.data = data
        throw err
    }

    return data
}
