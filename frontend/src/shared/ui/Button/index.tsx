import styles from './button.module.css'

interface Props {
  variant?: 'primary' | 'secondary'
  children: React.ReactNode
  onClick?: () => void
  className?: string
}

export default function Button({ variant = 'primary', children, onClick, className }: Props) {
  const base = variant === 'primary' ? styles.primary : styles.secondary
  return (
    <button
      className={className ? `${base} ${className}` : base}
      onClick={onClick}
    >
      {children}
    </button>
  )
}
