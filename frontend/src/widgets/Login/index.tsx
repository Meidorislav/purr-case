import { Widget } from '@xsolla/login-sdk'
import { useEffect, useRef } from 'react'

const widget = new Widget({
    projectId: import.meta.env.VITE_XSOLLA_LOGIN_ID,
    preferredLocale: 'en_US',
    clientId: import.meta.env.VITE_XSOLLA_CLIENT_ID,
    responseType: 'code',
    state: 'purr-case',
    redirectUri: import.meta.env.VITE_XSOLLA_RETURN_URL,
    scope: 'offline'
})


interface Props {
    onClose: () => void
}

export default function Login({ onClose }: Props) {
    const containerRef = useRef<HTMLDivElement>(null)

    useEffect(() => {
        if (containerRef.current) {
            widget.mount('xsolla-login-widget')
            widget.open()
        }

        return () => {
            widget.unmount()
        }
    }, [])

    return <div onClick={onClose}>
        <div
            id="xsolla-login-widget"
            ref={containerRef}
            onClick={e => e.stopPropagation()}
        />
    </div>
}