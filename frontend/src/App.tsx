import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Header from './widgets/Header'
import Home from './pages/Home'
import styles from './App.module.css'

export default function App() {
  return (
    <BrowserRouter>
      <div className={styles.app}>
        <Header />
        <main className={styles.main}>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/inventory" element={<div>Inventory</div>} />
            <Route path="/battlepass" element={<div>Battle Pass</div>} />
            <Route path="/cart" element={<div>Cart</div>} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  )
}
