import { useEffect, useState } from 'react'
import InventoryCard from '../InventoryCard'
import styles from './inventory-list.module.css'

interface InventoryItem {
  item_id: number
  name: string
  description: string
  image_url: string | null
}

// TODO: replace with actual token from auth when login is implemented
const TOKEN = 'eyJhbGciOiJSUzI1NiIsImtpZCI6ImM1MTI3M2M4LTRkZWQtNDMyYy1hODgyLWI5MjE0MGRhNDUwZCIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFsbWF6bGVuYTg2M0BnbWFpbC5jb20iLCJleHAiOjE3NzU2NDYzNDIsImdyb3VwcyI6W3siaWQiOjY2MjE3LCJuYW1lIjoiZGVmYXVsdCIsImlzX2FjdGl2ZSI6dHJ1ZSwiaXNfZGVmYXVsdCI6dHJ1ZX1dLCJpYXQiOjE3NzU1NTk5NDIsImlzX21hc3RlciI6dHJ1ZSwiaXNfbmV3Ijp0cnVlLCJpc3MiOiJodHRwczovL2xvZ2luLnhzb2xsYS5jb20iLCJqdGkiOiJkOTNkNDAyZS1jZjJkLTQxOWYtYThlMS01MTUxNjM0MWZhZTAiLCJwcm9qZWN0X2lkIjozMDQyMDksInB1Ymxpc2hlcl9pZCI6ODc3ODM1LCJzdWIiOiIwZWI0YTkwOC1hNDJkLTRlMTQtYWU1MC0wOGZhNGYxYzkyZWEiLCJ0eXBlIjoieHNvbGxhX2xvZ2luIiwidXNlcm5hbWUiOiJqYWxlbmNlcyIsInhzb2xsYV9sb2dpbl9hY2Nlc3Nfa2V5IjoiMXAtcG1ka1UwQU5rMWhKNzJGcGxfeS13R3Q4VzFsbGx1RHhqWW9WWlE4cyIsInhzb2xsYV9sb2dpbl9wcm9qZWN0X2lkIjoiMDU2MzlkOTUtNDM2ZS00MmY0LWIzODEtZTFmMWRhMGIzZDE1In0.ehdibugWkyN73c2lCp09dEkVBhsdhjrmS3dD6MQ6J-Pi0TfOVfz0bN-jWA2mhnR90QY44GO9sC0B_6Z-u4GgtXH3xPdmkk7Uo2Rqmv-PPnBCfc3PLw2jZkluzUaf7ZG05wEg_K4IvXDAePxX9HXBdw-guWXq2iigYVZiXI_ju5uYUA-ZrbNH_9jHLjY5Yu_sTNtU1lESv4VkAhMvo-Wo9jUVSehWfpC1GPlRbfZ6yM-C6FhBKGyRFfU44ZDFBaGsny9B-ViEcg7vLA2ggT-8ae0yBZjlSGQgQ0cS7n5ERyIQbBbtoVkoo13EYvxJXQbIS9qEkAw1wRrJ_YBwBUj6Hg'
 
export default function InventoryList() {
  const [items, setItems] = useState<InventoryItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetch('/api/inventory', {
      headers: { 'Authorization': `Bearer ${TOKEN}` },
    })
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch inventory')
        return res.json()
      })
      .then(data => setItems(Array.isArray(data) ? data : data.items ?? []))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  console.log(items)
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
