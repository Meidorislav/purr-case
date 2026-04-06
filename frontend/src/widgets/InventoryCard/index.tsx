import Button from '../../shared/ui/Button'
import styles from './inventory-card.module.css'

interface Props {
  image: string
  name: string
  description: string
  onOpen: () => void
}

export default function InventoryCard({ image, name, description, onOpen }: Props) {
  return (
    <div className={styles.card}>
      <div className={styles.imageWrapper}>
        <img src={image} alt={name} className={styles.image} />
      </div>
      <div className={styles.body}>
        <h3 className={styles.name}>{name}</h3>
        <p className={styles.description}>{description}</p>
        <Button variant="primary" className={styles.btn} onClick={onOpen}>Open</Button>
      </div>
    </div>
  )
}
