import { useCallback, useEffect, useState } from 'react'
import { useAuth } from '../../shared/hooks/useAuth'
import InventoryCard from '../InventoryCard'
import WonModal from '../WonModal'
import styles from './inventory-list.module.css'
import yarnIcon from '../../../assets/icons/yarn.png'
import fishIcon from '../../../assets/icons/fish.png'
import foodIcon from '../../../assets/icons/food.png'

interface InventoryItem {
  id: number
  sku: string
  name: string
  description: string
  image_url: string | null
  quantity: number
  type: string
  actions: string[] | null
  custom_attributes?: {
    rarity?: string
    type?: string
  }
}

interface WonItem {
  name: string
  description: string
  image_url: string | null
  custom_attributes?: { rarity?: string }
  groups: { external_id: string; name: string }[]
}

interface InventoryResponse {
  items: InventoryItem[]
  currencies: InventoryItem[]
}

const CURRENCY_STUBS: { sku: string; icon: string }[] = [
  { sku: 'yarn', icon: yarnIcon },
  { sku: 'fish', icon: fishIcon },
  { sku: 'food', icon: foodIcon },
]

export default function InventoryList() {
  const { fetchWithAuth } = useAuth()
  const [items, setItems] = useState<InventoryItem[]>([])
  const [currencies, setCurrencies] = useState<InventoryItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [wonItem, setWonItem] = useState<WonItem | null>(null)

  const loadInventory = useCallback(() => {
    setLoading(true)
    fetchWithAuth('/api/inventory')
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch inventory')
        return res.json() as Promise<InventoryResponse>
      })
      .then(data => {
        setItems(data.items ?? [])
        setCurrencies(data.currencies ?? [])
      })
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    loadInventory()
  }, [loadInventory])

  const handleAction = async (item: InventoryItem, action: string) => {
    const isOpen = action === 'open'
    const endpoint = isOpen ? `/api/cases/${item.sku}/open` : '/api/inventory/unpack'
    const body = isOpen ? undefined : JSON.stringify({ sku: item.sku, quantity: 1 })

    const res = await fetchWithAuth(endpoint, {
      method: 'POST',
      headers: body ? { 'Content-Type': 'application/json' } : undefined,
      body,
    })

    if (!res.ok) {
      const err = await res.json()
      alert(err.error ?? 'Something went wrong')
      return
    }

    if (isOpen) {
      const data = await res.json()
      if (data.won_item) {
        setWonItem(data.won_item)
      }
    }
    loadInventory()
  }
  
  return (
    <section>
      {wonItem && (
        <WonModal
          name={wonItem.name}
          description={wonItem.description}
          image_url={wonItem.image_url}
          rarity={wonItem.custom_attributes?.rarity ?? 'common'}
          groups={wonItem.groups ?? []}
          onClose={() => setWonItem(null)}
        />
      )}
      <div className={styles.hero}>
        <h1 className={styles.heroTitle}>Meow collection</h1>
        <p className={styles.heroSubtitle}>All the rare stuff your cat dragged in</p>
      </div>
      <div className={styles.header}>
        <h2 className={styles.title}>Inventory</h2>
        <div className={styles.gameCurrency}>
          {CURRENCY_STUBS.map(({ sku, icon }) => {
            const found = currencies.find(c => c.sku === sku)
            return (
              <span key={sku} className={styles.label}>
                <img className={styles.icon} src={found?.image_url ?? icon} alt={sku} />
                <p className={styles.value}>x {found?.quantity ?? 0}</p>
              </span>
            )
          })}
        </div>
      </div>
      {loading && items.length === 0 && <p className={styles.loading}>Loading...</p>}
      {error && <p>{error}</p>}
      <div className={styles.list}>
        {items.map(item => (
          <InventoryCard
            key={item.id}
            image={item.image_url ?? ''}
            name={item.name}
            description={item.description}
            quantity={item.quantity}
            rarity={item.custom_attributes?.rarity}
            actions={item.actions ?? []}
            onAction={(action) => handleAction(item, action)}
          />
        ))}
      </div>
    </section>
  )
}
