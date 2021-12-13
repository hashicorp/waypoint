import { Response } from 'miragejs';

export default interface RouteHandler {
  serialize(response: unknown, serializerType: string): Response;
}
