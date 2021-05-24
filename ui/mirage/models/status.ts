import { Model, belongsTo } from 'ember-cli-mirage';
import { Status } from 'waypoint-pb';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';

export default Model.extend({
  owner: belongsTo({ polymorphic: true }),

  toProtobuf(): Status {
    let result = new Status();

    result.setCompleteTime(dateToTimestamp(this.completeTime));
    result.setDetails(this.details);
    // result.setError
    // result.setExtension
    result.setStartTime(dateToTimestamp(this.startTime));
    result.setState(Status.State[this.state as keyof typeof Status.State]);

    return result;
  },
});

function dateToTimestamp(date: Date): Timestamp {
  let result = new Timestamp();

  result.setSeconds(Math.floor(date.valueOf() / 1000));

  return result;
}
