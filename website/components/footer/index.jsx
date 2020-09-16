import Link from 'next/link'

export default function Footer() {
  return (
    <footer className="g-footer">
      <div className="g-container">
        <div className="left">
          <Link href="/docs">
            <a>Docs</a>
          </Link>
        </div>
      </div>
    </footer>
  )
}
