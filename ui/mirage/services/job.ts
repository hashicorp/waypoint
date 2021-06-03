import { Response } from 'miragejs';
import { GetJobStreamResponse } from 'waypoint-pb';

export function stream(): Response {
  let result = new GetJobStreamResponse();

  // TODO(jgwhite): Implement GetJobStream handler

  return this.serialize(result, 'application');
}
