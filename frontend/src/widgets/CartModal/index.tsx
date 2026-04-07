import Button from '../../shared/ui/Button'
import { useCart } from '../../shared/hooks/useCart'
import type { CartItem as CartItemType } from '../../shared/hooks/useCart'
import CartItem from './CartItem'
import styles from './cart-modal.module.css'

interface Props {
  onClose: () => void
}

export default function CartModal({ onClose }: Props) {
  const { items, updateQuantity, clear, totalPrice } = useCart()

  // TODO: replace with actual token from auth when login is implemented
  const TOKEN = 'eyJhbGciOiJSUzI1NiIsImtpZCI6ImM1MTI3M2M4LTRkZWQtNDMyYy1hODgyLWI5MjE0MGRhNDUwZCIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFsbWF6bGVuYTg2M0BnbWFpbC5jb20iLCJleHAiOjE3NzU2NDYzNDIsImdyb3VwcyI6W3siaWQiOjY2MjE3LCJuYW1lIjoiZGVmYXVsdCIsImlzX2FjdGl2ZSI6dHJ1ZSwiaXNfZGVmYXVsdCI6dHJ1ZX1dLCJpYXQiOjE3NzU1NTk5NDIsImlzX21hc3RlciI6dHJ1ZSwiaXNfbmV3Ijp0cnVlLCJpc3MiOiJodHRwczovL2xvZ2luLnhzb2xsYS5jb20iLCJqdGkiOiJkOTNkNDAyZS1jZjJkLTQxOWYtYThlMS01MTUxNjM0MWZhZTAiLCJwcm9qZWN0X2lkIjozMDQyMDksInB1Ymxpc2hlcl9pZCI6ODc3ODM1LCJzdWIiOiIwZWI0YTkwOC1hNDJkLTRlMTQtYWU1MC0wOGZhNGYxYzkyZWEiLCJ0eXBlIjoieHNvbGxhX2xvZ2luIiwidXNlcm5hbWUiOiJqYWxlbmNlcyIsInhzb2xsYV9sb2dpbl9hY2Nlc3Nfa2V5IjoiMXAtcG1ka1UwQU5rMWhKNzJGcGxfeS13R3Q4VzFsbGx1RHhqWW9WWlE4cyIsInhzb2xsYV9sb2dpbl9wcm9qZWN0X2lkIjoiMDU2MzlkOTUtNDM2ZS00MmY0LWIzODEtZTFmMWRhMGIzZDE1In0.ehdibugWkyN73c2lCp09dEkVBhsdhjrmS3dD6MQ6J-Pi0TfOVfz0bN-jWA2mhnR90QY44GO9sC0B_6Z-u4GgtXH3xPdmkk7Uo2Rqmv-PPnBCfc3PLw2jZkluzUaf7ZG05wEg_K4IvXDAePxX9HXBdw-guWXq2iigYVZiXI_ju5uYUA-ZrbNH_9jHLjY5Yu_sTNtU1lESv4VkAhMvo-Wo9jUVSehWfpC1GPlRbfZ6yM-C6FhBKGyRFfU44ZDFBaGsny9B-ViEcg7vLA2ggT-8ae0yBZjlSGQgQ0cS7n5ERyIQbBbtoVkoo13EYvxJXQbIS9qEkAw1wRrJ_YBwBUj6Hg'

  const handleCheckout = async () => {
    try {
      const res = await fetch('/api/payments/checkout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${TOKEN}`,
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
              <Button variant="primary" onClick={handleCheckout}>Checkout</Button>
              <Button variant="secondary" onClick={clear}>Clear cart</Button>
              <p className={styles.checkoutNote}>* Login required to complete purchase</p>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
