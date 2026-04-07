import { CartItem as CartItemType } from '../../shared/hooks/useCart'
import styles from './cart-modal.module.css'

interface Props {
  item: CartItemType
  onIncrease: () => void
  onDecrease: () => void
}

export default function CartItem({ item, onIncrease, onDecrease }: Props) {
  return (
    <div className={styles.item}>
      <img src={item.image} alt={item.name} className={styles.image} />
      <div className={styles.info}>
        <span className={styles.name}>{item.name}</span>
        <span className={styles.price}>${(item.price * item.quantity).toFixed(2)}</span>
      </div>
      <div className={styles.qty}>
        <button className={styles.qtyBtn} onClick={onDecrease}>−</button>
        <span className={styles.qtyNum}>{item.quantity}</span>
        <button className={styles.qtyBtn} onClick={onIncrease}>+</button>
      </div>
    </div>
  )
}
