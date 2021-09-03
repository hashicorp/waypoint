import s from './style.module.css'

export default function Hero() {
  return (
    <header className={s.hero}>
      <h1 className={s.heroHeading}>
        Easy application deployment for <em>Kubernetes</em> and <em>ECS</em>
      </h1>
      <p className={s.heroDescription}>
        Waypoint is an application deployment tool for Kubernetes, ECS, and many
        other platforms. It allows developers to deploy, manage, and observe
        their applications through a consistent abstraction of the underlying
        infrastructure. Operators can deliver a PaaS experience for developers
        across multiple platforms and clouds.
      </p>
      <div className={s.heroGraphic}>
        <img src={require('./hero-graphic.png')} width={1685} alt="" />
      </div>
    </header>
  )
}
