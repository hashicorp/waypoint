import { useSession, getCsrfToken } from 'next-auth/client'
import LoadingIcon from './loading.svg?include'
import InlineSvg from '@hashicorp/react-inline-svg'
import styles from './auth-gate.module.css'
import Button from '@hashicorp/react-button'
import { useEffect, useState } from 'react'

export default function AuthGate({ children }) {
  const [session, loading] = useSession()

  if (loading)
    return <InlineSvg className={styles.loadingIconSpin} src={LoadingIcon} />
  return session ? (
    <>{children}</>
  ) : (
    <div className={styles.signInWrapper}>
      <SignInForm />
    </div>
  )
}

function SignInForm() {
  const [token, setToken] = useState(null)
  useEffect(() => {
    async function getToken() {
      const t = await getCsrfToken()
      if (t) {
        setToken(t)
      }
    }
    getToken()
  }, [token])
  return token ? (
    <section className={styles.signInWrapper}>
      <Form token={token} callbackUrl={window.location.href} />
    </section>
  ) : null
}

function Form({ callbackUrl, token }) {
  return (
    <form action="http://localhost:3000/api/auth/signin/okta" method="POST">
      <input type="hidden" name="csrfToken" value={token} />
      <input type="hidden" name="callbackUrl" value={callbackUrl} />
      <Button
        type="submit"
        title={`Sign in with Okta`}
        theme={{ variant: 'primary', background: 'dark' }}
      />
    </form>
  )
}
