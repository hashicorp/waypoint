import Service from '@ember/service';
import { WaypointClient } from 'waypoint-client';
import { Request, UnaryInterceptor, UnaryResponse } from 'grpc-web';

// The UnaryInterceptor interface is for the promise-based client.
class ExampleUnaryInterceptor implements UnaryInterceptor<any, any> {
  intercept(
    request: Request<any, any>,
    invoker: (request: Request<any, any>) => Promise<UnaryResponse<any, any>>
  ) {
    const reqMsg = request.getRequestMessage();
    console.log('request message: ', reqMsg);
    return invoker(request).then((response: UnaryResponse<any, any>) => {
      const responseMsg = response.getResponseMessage();
      console.log('response message: ', responseMsg);
      return response;
    });
  }
}

export default class ApiService extends Service {
  opts = { unaryInterceptors: [new ExampleUnaryInterceptor()] };
  client = new WaypointClient('http://localhost:1235', null, null);
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    api: ApiService;
  }
}
