import s from './style.module.css'
import Content from '@hashicorp/react-content'

export default function SecurityPage() {
  return (
    <div className={s.root}>
      <Content
        content={
          <>
            <h1>Security</h1>
            <p>
              We understand that many users place a high level of trust in
              HashiCorp and the tools we build. We apply best practices and
              focus on security to make sure we can maintain the trust of the
              community.
            </p>
            <p>
              We deeply appreciate any effort to disclose vulnerabilities
              responsibly.
            </p>
            <p>
              If you would like to report a vulnerability, please see the{' '}
              <a href="https://www.hashicorp.com/security">
                HashiCorp security page
              </a>
              , which has the proper email to communicate with as well as our
              PGP key.
            </p>
            <p>
              If you aren&apos;t reporting a security sensitive vulnerability,
              please open an issue on the standard{' '}
              <a href="https://github.com/hashicorp/waypoint">GitHub</a>{' '}
              repository.
            </p>
          </>
        }
      />
    </div>
  )
}
