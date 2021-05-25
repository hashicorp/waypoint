import styles from './WaypointDiagram.module.css'
import LogoList from 'components/logo-list'
import classNames from 'classnames'

export default function WaypointDiagram({ className }) {
  return (
    <div className={classNames(styles.waypointDiagram, className)}>
      <img
        className={styles.codeToDeployment}
        src={require('./img/code-to-deployment.svg')}
      />

      <div className={styles.details}>
        <div className={styles.detailItem}>
          <img className={styles.icon} src={require('./img/code.svg')} />
          <div>
            <h4>Your application code</h4>
            <LogoList
              className={styles.logoList}
              logos={[
                {
                  alt: 'Ruby',
                  url: require('./img/logos/ruby.svg'),
                },
                {
                  alt: 'React',
                  url: require('./img/logos/react.svg'),
                },
                {
                  alt: 'Node JS',
                  url: require('./img/logos/nodejs.svg'),
                },
                {
                  alt: 'Go',
                  url: require('./img/logos/go.svg'),
                },
                {
                  alt: 'Python',
                  url: require('./img/logos/python.svg'),
                },
                {
                  alt: 'Angular',
                  url: require('./img/logos/angular.svg'),
                },
                {
                  alt: 'NextJS',
                  url: require('./img/logos/nextjs.svg'),
                },
                {
                  alt: 'and More',
                  url: require('./img/logos/and-more.svg'),
                },
              ]}
            />
          </div>
        </div>

        <div className={styles.detailItem}>
          <img className={styles.icon} src={require('./img/platform.svg')} />
          <div>
            <h4>Your deployment platform</h4>
            <LogoList
              reverse
              className={styles.logoList}
              logos={[
                {
                  alt: 'Azure Container Service',
                  url: require('./img/logos/azure-container-service.svg'),
                },
                {
                  alt: 'Amazon ECS',
                  url: require('./img/logos/amazon-ecs.svg'),
                },
                {
                  alt: 'Nomad',
                  url: require('./img/logos/nomad.svg'),
                },
                {
                  alt: 'Kubernetes',
                  url: require('./img/logos/kubernetes.svg'),
                },

                {
                  alt: 'Google Cloud Run',
                  url: require('./img/logos/cloud-run.svg'),
                },
                {
                  alt: 'Docker',
                  url: require('./img/logos/docker.svg'),
                },
                {
                  alt: 'and More',
                  url: require('./img/logos/and-more.svg'),
                },
              ]}
            />
          </div>
        </div>
      </div>
    </div>
  )
}
