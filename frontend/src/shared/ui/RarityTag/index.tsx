import styles from './rarity-tag.module.css'

interface Props {
  rarity: string
  className?: string
}

export default function RarityTag({ rarity, className }: Props) {
  return (
    <span className={`${styles.tag} ${styles[rarity] ?? ''} ${className ?? ''}`}>{rarity.toUpperCase()}</span>
  )
}
