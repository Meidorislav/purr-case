import Catalog from '../../widgets/Catalog'
import logo from '../../../assets/logo.svg'
import styles from './home.module.css'

export default function Home() {
  return (
    <div>
      <div className={styles.hero}>
        <img src={logo} alt="logo" className={styles.heroLogo} />
        <p className={styles.heroText}>The cats choose you. You just pay for it.</p>
      </div>
      <div className={styles.catalog}>
        <Catalog />
      </div>
    </div>
  )
}
