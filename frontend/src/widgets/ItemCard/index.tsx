import { useState } from 'react'
import Button from '../../shared/ui/Button'
import ItemModal from '../ItemModal'
import styles from './item-card.module.css'

interface VirtualPrice {
  name: string
  amount: number
  image_url: string | null
}

interface Group {
  external_id: string
  name: string
}

interface ContentItem {
  item_id: number
  name: string
  image_url: string | null
  quantity: number
}

interface Props {
  image: string
  name: string
  description: string
  price: number
  rawPrice: { amount: string; currency: string } | null
  virtualPrices: VirtualPrice[]
  groups: Group[]
  canBeBought: boolean
  isFree: boolean
  content?: ContentItem[]
  addToCartClick: () => void
}

export default function ItemCard({
  image, name, description, price,
  rawPrice, virtualPrices, groups, canBeBought, isFree,
  content, addToCartClick,
}: Props) {
  const [modalOpen, setModalOpen] = useState(false)

  return (
    <>
      <div className={styles.card}>
        <div className={styles.imageWrapper}>
          <img src={image} alt={name} className={styles.image} />
        </div>
        <div className={styles.body}>
          <h3 className={styles.name}>{name}</h3>
          <p className={styles.description}>{description}</p>
          <p className={styles.price}>${price}</p>
          <Button variant="primary" className={styles.btn} onClick={addToCartClick}>Add to cart</Button>
          <Button variant="secondary" className={styles.btn} onClick={() => setModalOpen(true)}>View items</Button>
        </div>
      </div>

      {modalOpen && (
        <ItemModal
          name={name}
          description={description}
          image_url={image || null}
          price={rawPrice}
          virtual_prices={virtualPrices}
          groups={groups}
          can_be_bought={canBeBought}
          is_free={isFree}
          content={content}
          onAddToCart={addToCartClick}
          onClose={() => setModalOpen(false)}
        />
      )}
    </>
  )
}
