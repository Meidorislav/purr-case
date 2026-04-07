import { useState } from 'react'
import Button from '../../shared/ui/Button'
import QuantityControl from '../../shared/ui/QuantityControl'
import ItemModal from '../ItemModal'
import styles from './item-card.module.css'
import eventSvg from '../../../assets/event.svg'
import { useCart } from '../../shared/hooks/useCart'

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
  sku: string
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
  content, addToCartClick, isEvent, sku
}: Props) {
  const [modalOpen, setModalOpen] = useState(false)
  const { items, updateQuantity } = useCart()
  const inCart = items.some((i: { sku: string }) => i.sku === sku)
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
          {inCart
            ? <>   
            <div className={styles.inCartBtnSection}>
              <span className={styles.inCartBtn}>In cart</span>        
              <QuantityControl
                  quantity={items.find(i => i.sku === sku)!.quantity}
                  onIncrease={addToCartClick}
                  onDecrease={() => updateQuantity(sku, items.find(i => i.sku === sku)!.quantity - 1)}
                />
            </div> 

            </>

            : <Button variant="primary" className={styles.btn} onClick={addToCartClick}>Add to cart</Button>
          }

          <Button variant="secondary" className={styles.btn} onClick={() => setModalOpen(true)}>View items</Button>
        </div>
      </div>

      {modalOpen && (
        <ItemModal
          sku={sku}
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
