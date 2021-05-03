import Terminal from '@hashicorp/react-command-line-terminal'

/**
 * A FramedTerminal is a simple component which is responsible for determining
 * a given frame of a terminal if it were to be animated, given its specific
 * (zero based) frame.
 *
 * If a frame is passed that goes above the total number of frames, it will loop.
 * e.g. a FramedTerminal with 10 unique frames (0,1,2,3,4,5,6,7,8,9) will display
 * frame 0 if frame={10} is passed in.
 *
 * This FramedTerminal itself is purely responsible for determining the state
 * of the terminal at a given frame and passing down the props of what the given
 * frame would look like in the static <Terminal/>.
 *
 * As compared to the <Terminal />, each line has two notable distinctions:
 *   - The `code` prop can be an array, which will allow the animation of individual lines.
 *   - An additional `frames` prop, which specifies how many frames (default 1) each step in the line will take
 *      - If code is an array, each step will take this many frames
 *      - Frames can also be 0, which will cause the next element to render at the same time
 *
 * Example Usage w/ number [e.g. (0)] indicating which frame it appears:
 *
 *  <FramedTerminal
 *    frame={0}
 *    lines={[
 *      {
 *        frames: 1,
 *        code: [
 *          '» Building . . .',                         (0)
 *          '» Building . . . . . .',                   (1)
 *          '» Building . . . . . . . . . ',            (2)
 *          '» Building . . . . . . . . . . . .',       (3)
 *        ]
 *      },
 *      {
 *        frames: 0,
 *        color: 'gray',
 *        code: 'Creating new buildpack-based image:',  (4)
 *        indent: 1,
 *      },
 *      {
 *        frames: 2,                                    (4)
 *        color: 'gray',
 *        code: 'heroku/buildpacks:18',
 *        indent: 1,
 *      },
 *      {
 *        frames: 1,                                    (6)
 *        color: 'gray',
 *        code: 'Some more code',
 *        indent: 1,
 *      },
 *    ]}
 *  />
 */
export default function FramedTerminal({ frame, lines }) {
  // Determine the total number of frames
  let totalFrames = 0
  lines.forEach((line) => {
    let frames = line.frames ? line.frames : 1
    if (Array.isArray(line.code)) {
      totalFrames += line.code.length * frames
    } else {
      totalFrames += frames
    }
  })

  // Determine the actual activeFrame to handle when the number
  // of frames passed in exceeds our totalFrames
  const activeFrame = frame % totalFrames

  // Calculate the lines that should actively be displayed
  // and passed down to our terminal
  let previousFrames = 0
  const terminalLines = lines
    .map((line) => {
      // Determine how many frames left we have here
      var remainingFrames = activeFrame - previousFrames

      // Calculate our result for this line that will be passed down
      // to our <Terminal />
      let result = null
      if (remainingFrames >= 0) {
        if (!Array.isArray(line.code)) {
          result = {
            color: line.color,
            code: line.code,
            indent: line.indent,
            short: line.short,
          }
        } else {
          var lineFrame = Math.floor(remainingFrames / line.frames)
          result = {
            color: line.color,
            code: line.code.slice(0, lineFrame + 1).splice(-1, 1),
            indent: line.indent,
          }
        }
      }

      // Increment our previousFrames
      let lineFrames = line.frames ? line.frames : 1
      if (Array.isArray(line.code)) {
        previousFrames += line.code.length * lineFrames
      } else {
        previousFrames += lineFrames
      }
      return result
    })
    .filter((el) => el != null)

  return <Terminal lines={terminalLines} noScroll product="waypoint" />
}
