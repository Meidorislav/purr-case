import { NavLink, useLocation } from 'react-router-dom'
import styles from './Header.module.css'
import logo from '../../../assets/logo.svg'

interface NavItem {
  label: string
  to: string
}

interface Props {
  nav?: NavItem[]
}

const defaultNav: NavItem[] = [
  { label: 'Catalog', to: '/' },
  { label: 'Battle Pass', to: '/battlepass' },
  { label: 'Inventory', to: '/inventory' },
  { label: 'Cart', to: '/cart'},

]

export default function Header({ nav = defaultNav }: Props) {
  const { pathname } = useLocation()
  const isMain = pathname === '/'

  return (
    <header className={styles.header}>
      <div className={styles.inner}>
        {!isMain && (
          <div className={styles.logo}>
            <img src={logo} alt="logo" />
          </div>
        )}
        <nav className={styles.nav}>
          {nav.map(({ label, to }) => (
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
        </nav>
      </div>
    </header>
  )
}
