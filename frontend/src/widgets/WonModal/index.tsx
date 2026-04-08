import RarityTag from '../../shared/ui/RarityTag'
import styles from './WonModal.module.css'


interface Group {
  external_id: string
  name: string
}


interface Props {
  name: string
  description: string
  image_url: string | null
  onClose: () => void
  rarity: string
  groups: Group[]
}

export default function WonModal({
  name, description, image_url, onClose, rarity, groups
}: Props) {

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={`${styles.modal}`} onClick={e => e.stopPropagation()}>
        <button className={styles.closeBtn} onClick={onClose}>✕</button>
        <h1 className={styles.title}>A cat chose you!</h1>
        <div className={styles.body}>
            {image_url && (
              <div className={styles.imageWrapper}>
                <img src={image_url} alt={name} className={styles.image} />
              </div>
            )}

            <h2 className={styles.name}>{name}</h2>
            <p className={styles.description}>{description}</p>
            <div className={styles.tags}>
            {groups.map(g => (
                <span key={g.external_id} className={styles.tag}>{g.name}</span>
              ))}
              <RarityTag rarity={rarity} />
            </div>
      </div>
    </div>
  </div>
  )
}
