import { useState } from 'react'
import { NavLink, useLocation } from 'react-router-dom'
import styles from './Header.module.css'
import logo from '../../../assets/logo.svg'
import CartModal from '../CartModal'
import { useCart } from '../../shared/hooks/useCart'
import { useAuth } from '../../shared/hooks/useAuth'
import Login from '../Login'

const navLinks = [
  { label: 'Catalog', to: '/' },
  { label: 'Inventory', to: '/inventory' },
]

export default function Header() {
  const { pathname } = useLocation()
  const isMain = pathname === '/'
  const [cartOpen, setCartOpen] = useState(false)
  const [loginOpen, setLoginOpen] = useState(false)
  const { totalCount } = useCart()
  const { accessToken, logout } = useAuth()

  return (
    <>
      <header className={styles.header}>
        <div className={styles.inner}>
          {!isMain && (
            <div className={styles.logo}>
              <img src={logo} alt="logo" />
            </div>
          )}
          <nav className={styles.nav}>
            {navLinks.map(({ label, to }) => (
              <NavLink
                key={to}
                to={to}
                end
                className={({ isActive }) =>
                  isActive ? `${styles.link} ${styles.active}` : styles.link
                }
              >
                {label}
              </NavLink>
            ))}
            <button className={styles.link} onClick={() => setCartOpen(true)}>
              Cart{totalCount > 0 ? ` (${totalCount})` : ''}
            </button>
            {accessToken
              ? <button className={styles.link} onClick={logout}>Logout</button>
              : <button className={styles.link} onClick={() => setLoginOpen(true)}>Login</button>
            }
          </nav>
        </div>
      </header>

      {cartOpen && <CartModal onClose={() => setCartOpen(false)} onLoginRequest={() => { setCartOpen(false); setLoginOpen(true) }} />}
      {loginOpen && <Login onClose={() => setLoginOpen(false)} />}
    </>
  )
}
