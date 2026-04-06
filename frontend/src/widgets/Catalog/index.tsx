import { useEffect, useState } from 'react'
import ItemCard from '../ItemCard'
import styles from './catalog.module.css'

const FILTERS = ['All', 'Cases', 'Skins', 'Currency']

const FILTER_GROUPS: Record<string, string[]> = {
  Currency: ['Currency', 'Currency Packs'],
}

interface Group {
  external_id: string
  name: string
}

interface VirtualPrice {
  name: string
  amount: number
  image_url: string | null
}

interface ContentItem {
  item_id: number
  name: string
  image_url: string | null
  quantity: number
}

interface Item {
  item_id: number
  type: string
  name: string
  description: string
  image_url: string | null
  price: { amount: string; currency: string } | null
  virtual_prices: VirtualPrice[]
  groups: Group[]
  can_be_bought: boolean
  is_free: boolean
  content?: ContentItem[]
}

export default function Catalog() {
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeFilter, setActiveFilter] = useState('All')

  useEffect(() => {
    fetch('/api/items')
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch items')
        return res.json()
      })
      .then(data => setItems(data.items))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  const isEvent = (item: Item) => item.groups.some(g => g.name === 'Event')

  const filtered = (activeFilter === 'All'
    ? items
    : items.filter(item => {
        const groupNames = FILTER_GROUPS[activeFilter] ?? [activeFilter]
        return item.groups.some(g => groupNames.includes(g.name))
      })
  ).sort((a, b) => Number(isEvent(b)) - Number(isEvent(a)))
  return (
    <section className={styles.catalog}>
      <div className={styles.header}>
        <h2 className={styles.title}>CATalog</h2>
        <div className={styles.filters}>
          {FILTERS.map(f => (
            <button
              key={f}
              className={`${styles.filter} ${activeFilter === f ? styles.filterActive : ''}`}
              onClick={() => setActiveFilter(f)}
            >
              {f}
            </button>
          ))}
        </div>
      </div>
      {loading && <p className={styles.loading}>Loading...</p>}
      {error && <p>{error}</p>}
      <div className={styles.list}>
        {filtered.map(item => (
          <ItemCard
            key={item.item_id}
            image={item.image_url ?? ''}
            name={item.name}
            description={item.description}
            price={item.price ? parseFloat(item.price.amount) : 0}
            rawPrice={item.price}
            virtualPrices={item.virtual_prices ?? []}
            groups={item.groups ?? []}
            canBeBought={item.can_be_bought}
            isFree={item.is_free}
            content={item.content}
            addToCartClick={() => {}}
          />
        ))}
      </div>
    </section>
  )
}
