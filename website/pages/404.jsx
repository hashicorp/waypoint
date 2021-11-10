import Link from 'next/link'
import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'
import Placeholder from '@hashicorp/react-placeholder'

export default function NotFound() {
  const { asPath } = useRouter()
  const [isMounted, setIsMounted] = useState()

  useEffect(() => {
    setIsMounted(true)
  }, [])

  useEffect(() => {
    if (
      typeof window !== 'undefined' &&
      typeof window?.analytics?.track === 'function' &&
      typeof window?.document?.referrer === 'string' &&
      typeof window?.location?.href === 'string'
    )
      window.analytics.track(window.location.href, {
        category: '404 Response',
        label: window.document.referrer || 'No Referrer',
      })
  }, [])

  const defaultMessage = (
    <>
      <p>
        We&rsquo;re sorry but we can&rsquo;t find the page you&rsquo;re looking
        for.
      </p>
      <p>
        <Link href="/">
          <a>Back to Home</a>
        </Link>
      </p>
    </>
  )

  const reg = /\/(?<version>v\d+[.]\d+[.](\d+|x))/g
  const matches = reg.exec(asPath)
  const docsMessage = (
    <>
      <p>
        We&rsquo;re sorry, but this page does not exist for version&nbsp;
        <b>{matches?.groups?.version}</b>.
      </p>
      <p>
        Try viewing the&nbsp;
        <Link href={asPath.replace(reg, '')}>
          <a>latest version</a>
        </Link>
        &nbsp;instead.
      </p>
    </>
  )

  // If we're viewing a versioned page, provide a contextual message
  const errorMessage = matches ? docsMessage : defaultMessage

  return (
    <div id="p-404" className="g-grid-container">
      <h1 className="g-type-display-1">Page Not Found</h1>
      {isMounted ? (
        errorMessage
      ) : (
        <Placeholder>
          {(Box) => <Box lines={['60ch', '15ch']} prose />}
        </Placeholder>
      )}
    </div>
  )
}
