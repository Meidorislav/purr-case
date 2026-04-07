import Button from '../../shared/ui/Button'
import { useCart } from '../../shared/hooks/useCart'
import type { CartItem as CartItemType } from '../../shared/hooks/useCart'
import { useAuth } from '../../shared/hooks/useAuth'
import CartItem from './CartItem'
import styles from './cart-modal.module.css'

interface Props {
  onClose: () => void
  onLoginRequest: () => void
}

export default function CartModal({ onClose, onLoginRequest }: Props) {
  const { items, updateQuantity, clear, totalPrice } = useCart()
  const { accessToken, fetchWithAuth } = useAuth()

  const handleCheckout = async () => {
    if (!accessToken) {
      onLoginRequest()
      return
    }

    try {
      const res = await fetchWithAuth('/api/payments/checkout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          items: items.map((item: CartItemType) => ({ sku: item.sku, quantity: item.quantity })),
        }),
      })

      if (!res.ok) {
        const data = await res.json()
        alert(data.error ?? 'Checkout failed')
        return
      }

      const { checkoutUrl } = await res.json()
      window.location.href = checkoutUrl
    } catch {
      alert('Network error, please try again')
    }
  }

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={e => e.stopPropagation()}>
        <button className={styles.closeBtn} onClick={onClose}>✕</button>
        <h2 className={styles.title}>Cart</h2>

        {items.length === 0 ? (
          <p className={styles.empty}>Your cart is empty</p>
        ) : (
          <>
            <div className={styles.list}>
              {items.map((item: CartItemType) => (
                <CartItem
                  key={item.sku}
                  item={item}
                  onIncrease={() => updateQuantity(item.sku, item.quantity + 1)}
                  onDecrease={() => updateQuantity(item.sku, item.quantity - 1)}
                />
              ))}
            </div>

            <div className={styles.footer}>
              <div className={styles.total}>
                <span>Total</span>
                <span>${totalPrice.toFixed(2)}</span>
              </div>
              <Button variant="primary" onClick={handleCheckout}>
                {accessToken ? 'Checkout' : 'Login to Checkout'}
              </Button>
              <Button variant="secondary" onClick={clear}>Clear cart</Button>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
