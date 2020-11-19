import { useEffect, useState } from 'react'

// The API endpoint for status.hashicorp.com that allows us to check
// out the status of our individual services.
//
// See https://status.hashicorp.com/api for details
const statusEndpoint =
  'https://pdrzb3d64wsj.statuspage.io/api/v2/components.json'

// The service IDs that we're interested in reporting is there's an outage
const waypointServiceIDs = [
  'j0bgx9v6fcp9', // Waypoint
  'cpn6w07r5f2y', // Waypoint URL service
]

// A Hook that reports the status of our Waypoint services via status.hashicorp.com
// reports true if everything is OK, and false if something is being reported as
// non-operational.
export default function useWaypointServiceStatus() {
  const [statusOK, setStatusOK] = useState(true)
  useEffect(() => {
    fetch(statusEndpoint)
      .then((resp) => resp.json())
      .then((data) => {
        // Filter down to only the Waypoint services
        return data.components.filter((component) =>
          waypointServiceIDs.includes(component.id)
        )
      })
      .then((components) => {
        // Set the status to false if one of our Waypoint services is not operational
        if (components.some((component) => component.status != 'operational')) {
          setStatusOK(false)
        }
      })
  }, null)
  return statusOK
}
