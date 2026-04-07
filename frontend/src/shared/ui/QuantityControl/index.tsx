import styles from './quantity-control.module.css'

interface Props {
  quantity: number
  onIncrease: () => void
  onDecrease: () => void
}

export default function QuantityControl({ quantity, onIncrease, onDecrease }: Props) {
  return (
    <div className={styles.qty}>
      <button className={styles.btn} onClick={onDecrease}>−</button>
      <span className={styles.num}>{quantity}</span>
      <button className={styles.btn} onClick={onIncrease}>+</button>
    </div>
  )
}
