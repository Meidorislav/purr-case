import { useEffect } from "react"
import { useNavigate } from "react-router-dom"


export default function Callback() {
    const navigate = useNavigate()

    useEffect(() => {
        const searchParams = new URLSearchParams(window.location.search)
        const code = searchParams.get('code')
        if (!code) {
            navigate('/')
            return
        }

        const body = new URLSearchParams()
        body.set('grant_type', 'authorization_code')
        body.set('client_id', import.meta.env.VITE_XSOLLA_CLIENT_ID)
        body.set('code', code)
        body.set('redirect_uri', import.meta.env.VITE_XSOLLA_RETURN_URL)

        fetch('https://login.xsolla.com/api/oauth2/token', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            body: body.toString()
        }).then(res => res.json()).then(res => {
            localStorage.setItem('access_token', res.access_token)
            localStorage.setItem('refresh_token', res.refresh_token)
            navigate('/')
        }).catch(() => navigate('/'))
    },
        [])

    return <div>Logging in...</div>
}