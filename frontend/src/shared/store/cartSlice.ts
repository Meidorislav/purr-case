import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface CartItem {
  sku: string
  name: string
  image: string
  price: number
  quantity: number
}

const STORAGE_KEY = 'purrcase-cart'

function load(): CartItem[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : []
  } catch {
    return []
  }
}

interface CartState {
  items: CartItem[]
}

const initialState: CartState = {
  items: load(),
}

const cartSlice = createSlice({
  name: 'cart',
  initialState,
  reducers: {
    addItem(state, action: PayloadAction<Omit<CartItem, 'quantity'>>) {
      const existing = state.items.find(i => i.sku === action.payload.sku)
      if (existing) {
        existing.quantity += 1
      } else {
        state.items.push({ ...action.payload, quantity: 1 })
      }
      localStorage.setItem(STORAGE_KEY, JSON.stringify(state.items))
    },
    removeItem(state, action: PayloadAction<string>) {
      state.items = state.items.filter(i => i.sku !== action.payload)
      localStorage.setItem(STORAGE_KEY, JSON.stringify(state.items))
    },
    updateQuantity(state, action: PayloadAction<{ sku: string; quantity: number }>) {
      const { sku, quantity } = action.payload
      if (quantity <= 0) {
        state.items = state.items.filter(i => i.sku !== sku)
      } else {
        const item = state.items.find(i => i.sku === sku)
        if (item) item.quantity = quantity
      }
      localStorage.setItem(STORAGE_KEY, JSON.stringify(state.items))
    },
    clear(state) {
      state.items = []
      localStorage.removeItem(STORAGE_KEY)
    },
  },
})

export const { addItem, removeItem, updateQuantity, clear } = cartSlice.actions
export default cartSlice.reducer
