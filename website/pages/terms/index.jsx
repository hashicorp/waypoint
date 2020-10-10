import LongformPage from 'components/longform-page'

export default function TermsOfUsePage() {
  return (
    <LongformPage
      alert={
        <p>
          These Terms of Use apply to the Waypoint Url Service as hosted by
          HashiCorp. If you do not feel comfortable accepting the terms, you may
          disable the URL service or self-host the URL service.{' '}
          <a href="/docs/url">Learn more here</a>.
        </p>
      }
      title="Terms of Use"
    >
      <p>
        Please read these Terms of Use (“Agreement”) carefully before using the
        Services offered by HashiCorp, Inc. (“Company”). You may not use the
        Website or Services if you do not unconditionally agree to this
        Agreement. You represent that: (i) if you are an individual, you are of
        legal age to form a binding contract (and at least 13 years of age or
        older) or have a parent’s permission to do so; (ii) if you are entering
        into this Agreement on behalf of an entity (e.g., a corporation), you
        are an authorized to bind the entity to this Agreement and agree to this
        Agreement on the entity’s behalf; (iii) all registration information you
        submit is accurate and truthful and you will maintain its accuracy; and
        (iv) you are legally permitted to use and access the Services and take
        full responsibility for the selection and use of and access to the
        Services. This Agreement is void where prohibited by law, and the right
        to access the Services is revoked in such jurisdictions.
      </p>
      <p>
        <strong>1. ACCESS TO THE SERVICES.</strong>
        The waypointproject.io website and domain name and any other linked
        pages, features, content, or application services (the “Website”) and
        any related services (such as the Waypoint URL Service, as described
        more fully on the Website) offered by HashiCorp (the “Services”) are
        owned and operated by Company. Company may change, suspend or
        discontinue the Website or Services at any time, including the
        availability of any feature, database, or Content. Company may, in its
        sole discretion, modify this Agreement at any time by posting a notice
        on the Website, or by sending you a notice. Your use of the Services
        following such notification constitutes your acceptance of the terms and
        conditions of this Agreement as modified.
      </p>
      <p>
        <strong>2. SERVICES CONTENT.</strong> The Services and the “Content”
        (which includes, without limitation website URLs, applications or other
        information linked to by the website URLs, text, graphics, articles,
        photographs, images, and/or illustrations) may only be used in
        accordance with the terms of this Agreement. You warrant that you
        possess all rights necessary to provide such Content to Company and to
        grant Company the rights to use such information in connection with the
        Services.
      </p>
    </LongformPage>
  )
}
