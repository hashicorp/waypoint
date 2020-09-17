export default function Footer({ openConsentManager }) {
  return (
    <footer className="g-footer">
      <div className="g-container">
        <div className="left">
          <a onClick={openConsentManager}>Consent Manager</a>
        </div>
      </div>
    </footer>
  )
}
