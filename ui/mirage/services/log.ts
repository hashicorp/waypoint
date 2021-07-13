import { LogBatch } from 'waypoint-pb';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { Response } from 'miragejs';
import { getUnixTime } from 'date-fns';

export function stream(): Response {
  // TODO(jgwhite): Implement GetLogStream handler (+ models & factories)

  let result = new LogBatch();
  let entry = new LogBatch.Entry();
  let ts = new Timestamp();

  ts.setSeconds(getUnixTime(new Date(2021, 0, 1)));

  entry.setSource(LogBatch.Entry.Source.APP);
  entry.setTimestamp(ts);
  entry.setLine('Example log line');

  result.addLines(entry);

  return this.serialize(result, 'application');
}
