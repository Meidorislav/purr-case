import { useEffect, useState } from 'react'
import InventoryCard from '../InventoryCard'
import styles from './inventory-list.module.css'

interface InventoryItem {
  item_id: number
  name: string
  description: string
  image_url: string | null
}

export default function InventoryList() {
  const [items, setItems] = useState<InventoryItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // TODO: fetch inventory и отдельно валюту выводить
  useEffect(() => {
    fetch('/api/items')
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch inventory')
        return res.json()
      })
      .then(data => setItems(data.items))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  return (
    <section>
      <div className={styles.hero}>
        <h1 className={styles.heroTitle}>Meow collection</h1>
        <p className={styles.heroSubtitle}>All the rare stuff your cat dragged in</p>
      </div>
      <div className={styles.header}>
        <h2 className={styles.title}>Inventory </h2>
        <div className={styles.gameCurrency}>
          <span className={styles.label}>
            <img className={styles.icon} src="../../../assets/icons/yarn.png" alt="yarn" />
            <p className={styles.value}>x 300</p>
          </span>
          <span className={styles.label}>
            <img className={styles.icon} src="../../../assets/icons/fish.png" alt="fish" />
            <p className={styles.value}>x 100</p>
          </span>
          <span className={styles.label}>
            <img className={styles.icon} src="../../../assets/icons/food.png" alt="food" />
            <p className={styles.value}>x 200</p>
          </span>
        </div>
      </div>
      {loading && <p className={styles.loading}>Loading...</p>}
      {error && <p>{error}</p>}
      <div className={styles.list}>
        {items.map(item => (
          <InventoryCard
            key={item.item_id}
            image={item.image_url ?? ''}
            name={item.name}
            description={item.description}
            onOpen={() => {}}
          />
        ))}
      </div>
    </section>
  )
}
