import { useState } from 'react'
import Button from '../../shared/ui/Button'
import RarityTag from '../../shared/ui/RarityTag'
import styles from './inventory-card.module.css'

interface Props {
  image: string
  name: string
  description: string
  quantity: number
  rarity?: string
  actions: string[]
  onAction: (action: string) => Promise<void> | void
}

const ACTION_LABELS: Record<string, string> = {
  open: 'Open',
  unpack: 'Unpack',
}

const ACTION_LOADING_LABELS: Record<string, string> = {
  open: 'Opening...',
  unpack: 'Unpacking...',
}

export default function InventoryCard({ image, name, description, quantity, rarity, actions, onAction }: Props) {
  const [loadingAction, setLoadingAction] = useState<string | null>(null)

  const handleClick = async (action: string) => {
    setLoadingAction(action)
    try {
      await onAction(action)
    } finally {
      setLoadingAction(null)
    }
  }

  return (
    <div className={styles.card}>
      <div className={styles.imageWrapper}>
        <img src={image} alt={name} className={styles.image} />
        {quantity > 1 && <span className={styles.quantity}>x{quantity}</span>}
        {rarity && <RarityTag rarity={rarity} className={styles.rarity} />}
      </div>
      <div className={styles.body}>
        <h3 className={styles.name}>{name}</h3>
        <p className={styles.description}>{description}</p>
        {actions.map(action => (
          <Button key={action} variant="primary" className={styles.btn} disabled={!!loadingAction} onClick={() => handleClick(action)}>
            {loadingAction === action ? (ACTION_LOADING_LABELS[action] ?? action) : (ACTION_LABELS[action] ?? action)}
          </Button>
        ))}
      </div>
    </div>
  )
}
