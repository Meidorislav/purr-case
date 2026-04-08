import { useState } from 'react'
import InventoryList from '../../widgets/InventoryList'
import Login from '../../widgets/Login'
import { useAuth } from '../../shared/hooks/useAuth'
import styles from './inventory.module.css'

export default function Inventory() {
  const { accessToken } = useAuth()
  const [loginOpen, setLoginOpen] = useState(!accessToken)

  return (
    <div className={styles.page}>
      {loginOpen && <Login onClose={() => setLoginOpen(false)} />}
      <InventoryList />
    </div>
  )
}
