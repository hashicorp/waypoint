import { Fragment } from 'react'
import styles from './placement-table.module.css'

export default function PlacementTable({ groups = [] }) {
  return (
    <table className={styles.placementTable}>
      <thead>
        <tr>
          <td width="120" className={styles.placementHead}>
            Placement
          </td>
          <td className={styles.placementGroup}>
            {Array.isArray(groups[0]) ? (
              groups.map((subgroup) => {
                return (
                  <Fragment key={subgroup.join('')}>
                    <code
                      dangerouslySetInnerHTML={{
                        __html: wrapLastItem(subgroup, 'strong').join(' -> '),
                      }}
                    />
                    <br />
                  </Fragment>
                )
              })
            ) : (
              <code
                dangerouslySetInnerHTML={{
                  __html: wrapLastItem(groups, 'strong').join(' -> '),
                }}
              />
            )}
          </td>
        </tr>
      </thead>
    </table>
  )
}

function wrapLastItem(arr, wrapper) {
  arr[arr.length - 1] = `<${wrapper}>${arr[arr.length - 1]}</${wrapper}>`
  return arr
}
