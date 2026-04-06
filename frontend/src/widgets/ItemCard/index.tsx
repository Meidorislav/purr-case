import Button from '../../shared/ui/Button'
import styles from './item-card.module.css'

interface Props {
  image: string
  name: string
  description: string
  price: number
  addToCartClick: () => void
  viewItemsClick: () => void
}

export default function ItemCard({ image, name, description, price }: Props) {
  return (
    <div className={styles.card}>
      <div className={styles.imageWrapper}>
        <img src={image} alt={name} className={styles.image} />
      </div>
      <div className={styles.body}>
        <h3 className={styles.name}>{name}</h3>
        <p className={styles.description}>{description}</p>
        <p className={styles.price}>${price}</p>
        <Button variant="primary" className={styles.btn} onClick={() => {}}>Add to cart</Button>
        <Button variant="secondary" className={styles.btn} onClick={() => {}}>View items</Button>
      </div>
    </div>
  )
}
