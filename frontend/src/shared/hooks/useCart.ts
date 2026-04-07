import { useState, useEffect } from 'react'

export interface CartItem {
  sku: string
  name: string
  image: string
  price: number
  quantity: number

}

const STORAGE_KEY = 'purr-cart'

function load(): CartItem[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : []
  } catch {
    return []
  }
}

function save(items: CartItem[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(items))
}

export function useCart() {
  const [items, setItems] = useState<CartItem[]>(load)

  useEffect(() => {
    save(items)
  }, [items])

  const addItem = (item: Omit<CartItem, 'quantity'>) => {
    setItems(prev => {
      const existing = prev.find(i => i.sku === item.sku)
      if (existing) {
        return prev.map(i => i.sku === item.sku ? { ...i, quantity: i.quantity + 1 } : i)
      }
      return [...prev, { ...item, quantity: 1 }]
    })
  }

  const removeItem = (sku: string) => {
    setItems(prev => prev.filter(i => i.sku !== sku))
  }

  const updateQuantity = (sku: string, quantity: number) => {
    if (quantity <= 0) {
      removeItem(sku)
      return
    }
    setItems(prev => prev.map(i => i.sku === sku ? { ...i, quantity } : i))
  }

  const clear = () => setItems([])

  const totalCount = items.reduce((sum, i) => sum + i.quantity, 0)
  const totalPrice = items.reduce((sum, i) => sum + i.price * i.quantity, 0)

  return { items, addItem, removeItem, updateQuantity, clear, totalCount, totalPrice }
}
