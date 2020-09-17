import styles from './signin-error.module.css'
import Button from '@hashicorp/react-button'

export default function SigninErrorPage() {
  return (
    <div className={styles.signinErrorWrapper}>
      <h1 className="g-type-display-3">
        Sorry! <br />
        It seems you do not have appropriate permissions to view this content.
      </h1>
      <Button url="/" title="Go back" />
    </div>
  )
}
