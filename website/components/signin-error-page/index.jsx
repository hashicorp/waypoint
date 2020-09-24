import styles from './signin-error.module.css'
import Button from '@hashicorp/react-button'
import { useAuthProviders } from 'components/auth-gate'

export default function SigninErrorPage() {
  const authProviders = useAuthProviders()
  return (
    <div className={styles.signinErrorWrapper}>
      <h1 className="g-type-display-3">
        Sorry! <br />
        It seems you do not have appropriate permissions to view this content.
      </h1>
      {authProviders && (
        <div className={styles.logoutCard}>
          <div className={styles.authProviderGoBack}>
            <h4>{`If you'd like to try again with another Okta account, please log out of Okta`}</h4>
            {process.env.NEXT_PUBLIC_OKTA_DOMAIN && (
              <Button
                url={`https://${process.env.NEXT_PUBLIC_OKTA_DOMAIN}`}
                title={`Go to Okta`}
              />
            )}
          </div>
          <div className={styles.authProviderGoBack}>
            <h4>{`If you'd like to try again with another Auth0 account, please go back`}</h4>
            <Button url="/" title="Go back" />
          </div>
          )
        </div>
      )}
    </div>
  )
}
