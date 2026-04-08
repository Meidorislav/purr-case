import { BrowserRouter, Routes, Route, useSearchParams } from 'react-router-dom'
import { useEffect } from 'react'
import Header from './widgets/Header'
import Home from './pages/Home'
import Inventory from './pages/Inventory'
import styles from './App.module.css'
import Callback from './pages/Callback'
import { useCart } from './shared/hooks/useCart'

function PurchaseHandler() {
  const [params, setParams] = useSearchParams()
  const { clear } = useCart()

  useEffect(() => {
    if (params.get('status') === 'done') {
      clear()
      setParams({}, { replace: true })
    }
  }, [])

  return null
}

export default function App() {
  return (
    <BrowserRouter>
      <div className={styles.app}>
        <PurchaseHandler />
        <Header />
        <main className={styles.main}>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/inventory" element={<Inventory />} />
            <Route path="/callback" element={<Callback />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  )
}
