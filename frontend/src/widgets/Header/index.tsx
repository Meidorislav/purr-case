import { useState } from 'react'
import { NavLink, useLocation } from 'react-router-dom'
import styles from './Header.module.css'
import logo from '../../../assets/logo.svg'

const navLinks = [
  { label: 'Catalog', to: '/' },
  { label: 'Battle Pass', to: '/battlepass' },
  { label: 'Inventory', to: '/inventory' },
]

export default function Header() {
  const { pathname } = useLocation()
  const isMain = pathname === '/'
  const [cartOpen, setCartOpen] = useState(false)
  const [loginOpen, setLoginOpen] = useState(false)

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
            <button className={styles.link} onClick={() => setCartOpen(true)}>Cart</button>
            <button className={styles.link} onClick={() => setLoginOpen(true)}>Login</button>
          </nav>
        </div>
      </header>

      {cartOpen && (
        <div onClick={() => setCartOpen(false)}>
          {/* TODO: CartModal */}
        </div>
      )}
      {loginOpen && (
        <div onClick={() => setLoginOpen(false)}>
          {/* TODO: LoginModal */}
        </div>
      )}
    </>
  )
}
