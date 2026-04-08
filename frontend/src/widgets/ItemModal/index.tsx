import Button from '../../shared/ui/Button'
import QuantityControl from '../../shared/ui/QuantityControl'
import RarityTag from '../../shared/ui/RarityTag'
import { useCart } from '../../shared/hooks/useCart'
import styles from './item-modal.module.css'

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
  name: string
  description: string
  image_url: string | null
  price: { amount: string; currency: string } | null
  virtual_prices: VirtualPrice[]
  groups: Group[]
  can_be_bought: boolean
  is_free: boolean
  content?: ContentItem[]
  rarity?: string
  onAddToCart: () => void
  onClose: () => void
}

export default function ItemModal({
  sku, name, description, image_url, price,
  virtual_prices, groups, can_be_bought, is_free,
  content, rarity, onAddToCart, onClose,
}: Props) {
  const hasContent = content && content.length > 0
  const { items, updateQuantity } = useCart()
  const cartItem = items.find(i => i.sku === sku)
  const inCart = !!cartItem

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={`${styles.modal} ${hasContent ? styles.wide : ''}`} onClick={e => e.stopPropagation()}>
        <button className={styles.closeBtn} onClick={onClose}>✕</button>

        <div className={styles.body}>
          <div className={styles.left}>
            {image_url && (
              <div className={styles.imageWrapper}>
                <img src={image_url} alt={name} className={styles.image} />
              </div>
            )}

            <h2 className={styles.title}>{name}</h2>
            <p className={styles.description}>{description}</p>

            {(rarity || groups.length > 0) && (
              <div className={styles.groups}>
                {rarity && <RarityTag  className={styles.rarity} rarity={rarity} />}
                {groups.map(g => (
                  <span key={g.external_id} className={styles.tag}>{g.name}</span>
                ))}
              </div>
            )}

            <div className={styles.pricing}>
              {is_free && <span className={styles.free}>Free</span>}
              {price && (
                <span className={styles.price}>{price.amount} {price.currency}</span>
              )}
              {virtual_prices.length > 0 && (
                <div className={styles.virtualPrices}>
                  {virtual_prices.map((vp, i) => (
                    <div key={i} className={styles.virtualPrice}>
                      {vp.image_url && <img src={vp.image_url} alt={vp.name} className={styles.vpIcon} />}
                      <span>{vp.amount} {vp.name}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {hasContent && (
            <div className={styles.right}>
              <h3 className={styles.sectionTitle}>Items in case</h3>
              <div className={styles.contentGrid}>
                {content!.map(item => (
                  <div key={item.item_id} className={styles.contentItem}>
                    {item.image_url && (
                      <img src={item.image_url} alt={item.name} className={styles.contentImage} />
                    )}
                    <p className={styles.contentName}>{item.name}</p>
                    {item.quantity > 1 && (
                      <p className={styles.contentQty}>x{item.quantity}</p>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        {can_be_bought && (
          <div className={styles.footer}>
            {inCart
              ? <div className={styles.inCartBtnSection}>
                  <span className={styles.inCartBtn}>In cart</span>
                  <QuantityControl
                    quantity={cartItem!.quantity}
                    onIncrease={onAddToCart}
                    onDecrease={() => updateQuantity(sku, cartItem!.quantity - 1)}
                  />
                </div>
              : <Button variant="primary" className={styles.addBtn} onClick={onAddToCart}>Add to cart</Button>
            }
          </div>
        )}
      </div>
    </div>
  )
}
