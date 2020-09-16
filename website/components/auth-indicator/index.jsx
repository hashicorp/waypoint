import { signIn, signOut, useSession } from 'next-auth/client'
import LoadingIcon from './loading.svg?include'
import InlineSvg from '@hashicorp/react-inline-svg'
import styles from './auth-indicator.module.css'
import Button from '@hashicorp/react-button'

export default function AuthIndicator() {
  const [session, loading] = useSession()
  if (loading)
    return <InlineSvg className={styles.loadingIconSpin} src={LoadingIcon} />
  return (
    <div className={styles.authIndicator}>
      {session ? (
        <>
          <span className={`g-type-label ${styles.loggedInText}`}>
            Signed in as {session.user.email}
          </span>
          <span>
            <Button
              onClick={signOut}
              title={`Sign out`}
              size="small"
              theme={{ variant: 'secondary-neutral', background: 'dark' }}
            />
          </span>
        </>
      ) : (
        <Button
          onClick={signIn}
          title={`Sign in with Okta`}
          size="small"
          theme={{ variant: 'secondary-neutral', background: 'dark' }}
        />
      )}
    </div>
  )
}
