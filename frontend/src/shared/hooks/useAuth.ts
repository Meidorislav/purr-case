import { useState, useEffect } from 'react'

export function useAuth() {
    const [accessToken, setAccessToken] = useState<string | null>(localStorage.getItem('access_token'))

    async function refresh(): Promise<string | null> {
        const refreshToken = localStorage.getItem('refresh_token')
        if (!refreshToken) return null

        const body = new URLSearchParams()
        body.set('grant_type', 'refresh_token')
        body.set('refresh_token', refreshToken)
        body.set('client_id', import.meta.env.VITE_XSOLLA_CLIENT_ID)

        try {
            const res = await fetch('https://login.xsolla.com/api/oauth2/token', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                },
                body: body.toString()
            })

            if (!res.ok) {
                logout()
                return null
            }

            const data = await res.json()
            setAccessToken(data.access_token)
            localStorage.setItem('access_token', data.access_token)
            localStorage.setItem('refresh_token', data.refresh_token)
            return data.access_token
        } catch {
            logout()
            return null
        }
    }

    function logout() {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        setAccessToken(null)
    }

    async function fetchWithAuth(url: string, options: RequestInit = {}) {
        const token = localStorage.getItem('access_token')

        const res = await fetch(
            url,
            {
                ...options,
                headers: {
                    ...options.headers,
                    'Authorization': `Bearer ${token}`
                }
            }
        )

        if (res.status === 401) {
            const newToken = await refresh()
            if (!newToken) return res

            return fetch(url, {
                ...options,
                headers: {
                    ...options.headers,
                    'Authorization': `Bearer ${newToken}`
                }
            })
        }

        return res
    }

    return { accessToken, logout, fetchWithAuth }
}