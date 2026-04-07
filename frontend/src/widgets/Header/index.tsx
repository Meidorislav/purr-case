import { useState } from 'react'
import { NavLink, useLocation } from 'react-router-dom'
import styles from './Header.module.css'
import logo from '../../../assets/logo.svg'
import CartModal from '../CartModal'
import { useCart } from '../../shared/hooks/useCart'

const navLinks = [
  { label: 'Catalog', to: '/' },
  { label: 'Inventory', to: '/inventory' },
]

export default function Header() {
  const { pathname } = useLocation()
  const isMain = pathname === '/'
  const [cartOpen, setCartOpen] = useState(false)
  // const [loginOpen, setLoginOpen] = useState(false)
  const { totalCount } = useCart()

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
            <button className={styles.link}>Login</button>
          </nav>
        </div>
      </header>

      {cartOpen && <CartModal onClose={() => setCartOpen(false)} />}
      {/* {loginOpen && <LoginModal onClose={() => setLoginOpen(false)} />} */}
    </>
  )
}
