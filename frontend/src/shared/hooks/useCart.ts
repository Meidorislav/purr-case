import { useSelector, useDispatch } from 'react-redux'
import type { RootState, AppDispatch } from '../../app/store'
import { addItem, removeItem, updateQuantity, clear } from '../store/cartSlice'
import type { CartItem } from '../store/cartSlice'

export type { CartItem }

export function useCart() {
  const dispatch = useDispatch<AppDispatch>()
  const items = useSelector((state: RootState) => state.cart.items)

  const totalCount = items.reduce((sum: number, i: CartItem) => sum + i.quantity, 0)
  const totalPrice = items.reduce((sum: number, i: CartItem) => sum + i.price * i.quantity, 0)

  return {
    items,
    addItem: (item: Parameters<typeof addItem>[0]) => dispatch(addItem(item)),
    removeItem: (sku: string) => dispatch(removeItem(sku)),
    updateQuantity: (sku: string, quantity: number) => dispatch(updateQuantity({ sku, quantity })),
    clear: () => dispatch(clear()),
    totalCount,
    totalPrice,
  }
}

export default useCart
