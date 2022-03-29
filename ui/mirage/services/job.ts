import { Response } from 'miragejs';
import { RouteHandler } from '../types';
import { GetJobStreamResponse } from 'waypoint-pb';
import { dateToTimestamp } from '../utils';

export function stream(this: RouteHandler): Response {
  let result = new GetJobStreamResponse();
  let terminal = new GetJobStreamResponse.Terminal();

  terminal.setEventsList([
    event('Deploying v1...'),
    event('Deploying with Kubernetes...'),
    event('Creating resources...'),
    event('Having some cake...'),
    event('This cake is delicious...'),
    event('Would you like some?'),
    event('OK, suit yourself'),
  ]);

  result.setTerminal(terminal);

  // TODO(jgwhite): Implement GetJobStream handler

  return this.serialize(result, 'application');
}

function event(msg) {
  let result = new GetJobStreamResponse.Terminal.Event();
  let line = new GetJobStreamResponse.Terminal.Event.Line();

  line.setMsg(msg);

  result.setLine(line);
  result.setTimestamp(dateToTimestamp(new Date()));

  return result;
}
