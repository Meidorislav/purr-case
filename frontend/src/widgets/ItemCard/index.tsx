import { useState } from 'react'
import Button from '../../shared/ui/Button'
import ItemModal from '../ItemModal'
import styles from './item-card.module.css'
import eventSvg from '../../../assets/event.svg'

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
  isEvent?: boolean
}

export default function ItemCard({
  image, name, description, price,
  rawPrice, virtualPrices, groups, canBeBought, isFree,
  content, addToCartClick, isEvent
}: Props) {
  const [modalOpen, setModalOpen] = useState(false)
  isEvent = groups?.map(c => c.name).includes('Event')
  return (
    <>
      <div className={`${styles.card} ${isEvent ? styles.cardEvent : ''}`}>
        {isEvent && <img src={eventSvg} className={styles.event} alt="event" />}
        <div className={styles.imageWrapper}>
          <img src={image} alt={name} className={styles.image} />
        </div>
        <div className={styles.body}>
          <h3 className={styles.name}>{name}</h3>
          <p className={styles.description}>{description}</p>
          <div>
            <p className={styles.price}>${price}</p>
          </div>
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
