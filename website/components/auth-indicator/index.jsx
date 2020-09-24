import { signOut, useSession } from 'next-auth/client'
import styles from './auth-indicator.module.css'
import Button from '@hashicorp/react-button'

export default function AuthIndicator() {
  const [session, loading] = useSession()
  if (loading) return `Loading...`
  return (
    <div className={styles.authIndicator}>
      {session && (
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
      )}
    </div>
  )
}
