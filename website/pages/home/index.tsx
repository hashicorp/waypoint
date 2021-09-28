import Hero from 'components/homepage/hero'
import SectionIntro from 'components/homepage/sections/intro'
import SectionHowItWorks from 'components/homepage/sections/how-it-works'
import ConfigureYourApp from 'components/homepage/sections/how-it-works/configure-your-app'
import BuildAndDeploy from 'components/homepage/sections/how-it-works/build-and-deploy'
import MonitorAndManage from 'components/homepage/sections/how-it-works/monitor-and-manage'
import SectionMonitorAppHealth from 'components/homepage/sections/monitor-app-health'
import SectionExtendPlugins from 'components/homepage/sections/extend-plugins'
// import SectionWorkflowThatScales from 'components/homepage/sections/workflow-that-scales'
import s from './style.module.css'

export default function HomePage(): JSX.Element {
  return (
    <div className={s.homePage}>
      <Hero
        heading={<>Get the PaaS experience on your platform</>}
        description="Waypoint is an application deployment tool for Kubernetes, ECS, and
        many other platforms. It allows developers to deploy, manage, and
        observe their applications through a consistent abstraction of the
        underlying infrastructure."
      />
      <SectionIntro
        columnLeft={{
          heading: (
            <>
              Simple <em>developer experience</em>
            </>
          ),
          description:
            'Waypoint provides a simple and consistent abstraction for developers to easily build, deploy, and release applications.',
          features: [
            {
              icon: require('components/homepage/icons/layers.svg?include'),
              heading: 'Application-centric abstraction',
              description:
                'Specify the deployment needs with a simple and consistent abstraction without the underlying complexity.',
            },
            {
              icon: require('components/homepage/icons/link.svg?include'),
              heading: 'End-to-end deployment workflow',
              description:
                'Move and manage resources efficiently with distinct build, deploy, release steps.',
              link: {
                text: 'Learn more',
                url: 'https://www.waypointproject.io/docs/lifecycle',
              },
            },
          ],
        }}
        columnRight={{
          heading: (
            <>
              Powerful for <em>operators</em>
            </>
          ),
          description:
            'Waypoint enables operators to create PaaS workflows of Kubernetes, ECS, serverless applications.',
          features: [
            {
              icon: require('components/homepage/icons/maximize.svg?include'),
              heading: 'Build-deploy-release extensibility',
              description:
                'Enable a pluggable framework, integrated with CI/CD pipelines, monitoring tools, and any other ecosystem tools.',
              link: {
                text: 'Learn more',
                url: 'https://www.waypointproject.io/docs/extending-waypoint',
              },
            },
            {
              icon: require('components/homepage/icons/sidebar.svg?include'),
              heading: 'PaaS experience for developers',
              description:
                'Provide a consistent abstraction and unified workflow for any major platforms.',
            },
          ],
        }}
      />
      <SectionHowItWorks>
        <ConfigureYourApp
          heading="Configure your app for Waypoint"
          features={[
            {
              icon: require('components/homepage/icons/edit-pencil.svg?include'),
              heading: 'Writing waypoint.hcl files',
              description:
                'Your waypoint.hcl file defines how Waypoint builds, deploys, and releases a project.',
            },
            {
              icon: require('components/homepage/icons/layout.svg?include'),
              heading: 'Sample Waypoint files',
              description:
                'View sample waypoint.hcl files to see how straight-forward it is to configure your deployments',
            },
          ]}
          code={
            '<span class="token keyword">project =</span> "marketing-public"\n<span class="token keyword">app</span> "tech-blog" <span class="token keyword">{</span>\n<span class="token keyword">  build {</span>\n<span class="token keyword">    use</span> "pack" <span class="token keyword">{}</span> <span class="token comment"># Use Cloud Buildpacks</span>\n<span class="token keyword">  }</span>\n​\n<span class="token keyword">  deploy {</span>\n<span class="token keyword">    use</span> "kubernetes" <span class="token keyword">{}</span> <span class="token comment"># Deploy to Kubernetes</span>\n<span class="token keyword">  }</span>\n<span class="token keyword">}</span>'
          }
          codeNote="Configure your app for Waypoint in just a few lines"
        />
        <BuildAndDeploy
          heading="Build and deploy"
          features={[
            {
              icon: require('components/homepage/icons/file-plus.svg?include'),
              heading: 'Manage all steps within Waypoint',
              description:
                'Perform the build, deploy, and release steps for the app within waypoint. Or instrument your Waypoint deployments through Remote or Git operations.',
            },
          ]}
        />
        <MonitorAndManage
          heading="Manage your apps in one place"
          features={[
            {
              icon: require('components/homepage/icons/sliders.svg?include'),
              heading: 'Rich GUI for Waypoint',
              description:
                'No matter where your developers are deploying to, view logs, builds, releases and even run exec commands from the Waypoint UI.',
            },
          ]}
        />
      </SectionHowItWorks>
      <SectionMonitorAppHealth
        heading="Monitor app health on any cloud"
        description="One place to monitor the entire lifecycle of your applications, no
          matter where you deploy to. View Logs, Builds, Releasese and even run
          Exec commands from the Waypoint UI"
        features={[
          {
            icon: require('components/homepage/icons/eye.svg?include'),
            heading: 'A single pane of glass',
            description:
              'View all deployments, regardless of target from one location',
          },
        ]}
      />
      <SectionExtendPlugins
        heading="Extend Waypoint with plugins"
        description="Extend workflows via built-in plugins and an extensible interface. Supports custom builders, deployment platforms, registries, release managers, and more."
        features={[
          {
            icon: require('components/homepage/icons/box.svg?include'),
            heading: 'Available plugins',
            description: 'View a list of existing HashiCorp maintained plugins',
            link: {
              url: '/plugins',
              text: 'Plugins',
            },
          },
          {
            icon: require('components/homepage/icons/code-union.svg?include'),
            heading: 'Creating Waypoint plugins',
            description: 'Learn to extend Waypoint for your project’s needs',
            link: {
              url: '/docs/extending-waypoint/creating-plugins',
              text: 'Create',
            },
          },
        ]}
      />
      {/* <SectionWorkflowThatScales /> */}
    </div>
  )
}
