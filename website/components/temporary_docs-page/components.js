// -----------------------------------------------------
//                This code is LOCKED
//
// If any changes are needed to this code, or if this code
// is needed in any other projects, instead of changing or
// using it, instead we must complete this task as a prerequisite
//
// https://app.asana.com/0/1100423001970639/1195001770724993
//
// ------------------------------------------------------

import defaultMdxComponents from '@hashicorp/nextjs-scripts/lib/providers/docs'

// TODO: Jeff needs to make this reusable in some form and goes against his
// code lock above. I'm very sorry but needed this in the short term.
import Placement from '../placement-table'

export default function generateComponents(productName) {
  return defaultMdxComponents({
    product: productName,
    additionalComponents: { Placement },
  })
}
