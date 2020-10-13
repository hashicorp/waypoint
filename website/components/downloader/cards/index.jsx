import Tabs from '@hashicorp/react-tabs'
import {
  prettyOs,
  prettyArch,
  trackDownload,
} from 'components/downloader/utils/downloader'
import styles from './style.module.css'

export default function DownloadTabs({
  defaultTabIdx,
  tabData,
  downloads,
  version,
  merchandisingSlot,
  brand,
  logo,
  tutorialLink,
}) {
  return (
    <Tabs
      key={defaultTabIdx}
      centered
      fullWidthBorder
      theme={brand}
      className={styles.tabs}
      defaultTabIdx={defaultTabIdx}
      items={tabData.map(({ os, packageManagers }) => ({
        heading: prettyOs(os),
        tabChildren: function TabChildren() {
          return (
            <div className={styles.cards}>
              <Cards
                key={os}
                os={os}
                downloads={downloads}
                packageManagers={packageManagers}
                version={version}
                theme={brand}
                logo={logo}
                tutorialLink={tutorialLink}
              />
              {merchandisingSlot}
            </div>
          )
        },
      }))}
    />
  )
}

function Cards({
  os,
  downloads,
  packageManagers,
  version,
  theme,
  logo,
  tutorialLink,
}) {
  const arches = downloads[os]
  const hasPackageManager = Boolean(packageManagers)
  const hasMultiplePackageManagers = Array.isArray(packageManagers)

  function handleArchClick(arch) {
    return () => trackDownload('waypoint', version, os, arch)
  }

  return (
    <>
      <div
        className={
          hasMultiplePackageManagers
            ? styles.downloadCardsSingle
            : styles.downloadCards
        }
      >
        {hasPackageManager && (
          <div className={styles.packageManagers}>
            <span className={styles.cardTitle}>Package Manager</span>
            {Array.isArray(packageManagers) ? (
              <Tabs
                theme={theme}
                items={packageManagers.map(({ label, commands }) => ({
                  heading: label,
                  tabChildren: function TabChildren() {
                    return (
                      <div className={styles.install}>
                        {commands.map((command) => (
                          <pre key={command}>{command}</pre>
                        ))}
                      </div>
                    )
                  },
                }))}
              />
            ) : (
              <div className={styles.install}>
                {packageManagers.commands.map((command) => (
                  <pre key={command}>{command}</pre>
                ))}
              </div>
            )}
            {tutorialLink && (
              <div>
                <a href={tutorialLink.href}>{tutorialLink.label}</a>
              </div>
            )}
          </div>
        )}
        <div className={hasPackageManager ? styles.card : styles.soloCard}>
          <span className={styles.cardTitle}>Binary Download</span>
          <div className={styles.logoDownloadWrapper}>
            <div className={styles.logoWrapper}>
              {logo}
              <span className={styles.version}>{version}</span>
            </div>
            {Object.entries(arches).map(([arch, url]) => (
              <a
                href={url}
                key={arch}
                className={styles.downloadLink}
                onClick={handleArchClick(arch)}
              >
                {prettyArch(arch)}
              </a>
            ))}
          </div>
          <div className={styles.fastly}>
            Bandwidth courtesy of
            <img src={require('../logos/fastly.svg')} alt="Fastly" />
          </div>
        </div>
      </div>
    </>
  )
}
