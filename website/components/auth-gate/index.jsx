import { useSession, getCsrfToken, getProviders } from 'next-auth/client'
import LoadingIcon from './loading.svg?include'
import InlineSvg from '@hashicorp/react-inline-svg'
import styles from './auth-gate.module.css'
import Button from '@hashicorp/react-button'
import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'
import SigninErrorPage from 'components/signin-error-page'

export default function AuthGate({ children }) {
  const [session, loading] = useSession()
  const { query } = useRouter()
  if (query?.error === 'oAuthCallback') return <SigninErrorPage />

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
      <Form token={token} callbackUrl={window.location.origin} />
    </section>
  ) : null
}

function Form({ callbackUrl, token }) {
  const authProviders = useAuthProviders()

  return (
    <div className={styles.loginLockup}>
      {authProviders &&
        Object.keys(authProviders).map((ap) => (
          <form
            key={ap}
            action={`${callbackUrl}/api/auth/signin/${ap}`}
            method="POST"
          >
            <input type="hidden" name="csrfToken" value={token} />
            <input type="hidden" name="callbackUrl" value={callbackUrl} />
            <Button
              type="submit"
              title={`Sign in with ${authProviders[ap].name}`}
              theme={{ variant: 'primary', background: 'dark' }}
            />
          </form>
        ))}
    </div>
  )
}

export function useAuthProviders() {
  const [authProviders, setAuthProviders] = useState([])
  async function getProviderList() {
    const providers = await getProviders()
    if (providers) {
      setAuthProviders(providers)
    }
  }
  useEffect(() => {
    if (typeof window !== 'undefined') {
      getProviderList()
    }
  }, [])
  return authProviders
}
