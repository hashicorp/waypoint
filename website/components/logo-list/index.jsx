import styles from './LogoList.module.css'
import classNames from 'classnames'

export default function LogoList({ logos, className }) {
  return (
    <div className={classNames(styles.logoList, className)}>
      <ul>
        {logos.map((logo) => (
          <li key={logo.url}>
            <img src={logo.url} alt={logo.alt} />
          </li>
        ))}
      </ul>
    </div>
  )
}
