import React from 'react'
import styles from './Header.module.css'
import logo from '../../../assets/logo.svg'

export default function Header() {
  return (
    <header className={styles.header}>
      <div className={styles.inner}>
        <div className={styles.logo}>
          <img src={logo} alt="logo" />
        </div>
      </div>
    </header>
  )
}
