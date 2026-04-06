import styles from './case-modal.module.css'

interface CaseItem {
  item_id: number
  name: string
  image_url: string | null
  drop_chance?: number
}

interface Props {
  name: string
  description: string
  items: CaseItem[]
  onAddToCart: () => void
  onClose: () => void
}

function Stars({ value }: { value: number }) {
  return (
    <div className={styles.stars}>
      {[1, 2, 3].map(i => (
        <span key={i} className={i <= value ? styles.starFilled : styles.starEmpty}>★</span>
      ))}
    </div>
  )
}

export default function CaseModal({ name, description, items, onAddToCart, onClose }: Props) {
  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={e => e.stopPropagation()}>
        <h2 className={styles.title}>{name}</h2>
        <p className={styles.description}>{description}</p>

        <h3 className={styles.sectionTitle}>Items in case</h3>
        <div className={styles.carousel}>
          {items.map(item => (
            <div key={item.item_id} className={styles.item}>
              <img
                src={item.image_url ?? ''}
                alt={item.name}
                className={styles.itemImage}
              />
              <p className={styles.itemName}>{item.name}</p>
              <Stars value={2} />
              {item.drop_chance !== undefined && (
                <p className={styles.dropChance}>{item.drop_chance}%</p>
              )}
            </div>
          ))}
        </div>

        <button className={styles.addBtn} onClick={onAddToCart}>Add to cart</button>
      </div>
    </div>
  )
}
