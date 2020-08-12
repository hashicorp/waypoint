import Service from '@ember/service';
import { WaypointClient } from 'waypoint-client';
// import { Request, UnaryInterceptor, UnaryResponse } from 'grpc-web';

// The UnaryInterceptor interface is an interceptor example for the promise-based client that
// makes requests to the grpc-web proxy. Keeping this in source for now
// given that it will likely be useful for future feature.
//
// class ExampleUnaryInterceptor implements UnaryInterceptor<any, any> {
//   intercept(
//     request: Request<any, any>,
//     invoker: (request: Request<any, any>) => Promise<UnaryResponse<any, any>>
//   ) {
//     const reqMsg = request.getRequestMessage();
//     console.log('request message: ', reqMsg);
//     return invoker(request).then((response: UnaryResponse<any, any>) => {
//       const responseMsg = response.getResponseMessage();
//       console.log('response message: ', responseMsg);
//       return response;
//     });
//   }
// }
export default class ApiService extends Service {
  token = '';
  credentails = { authorization: this.token };
  // opts = { unaryInterceptors: [new ExampleUnaryInterceptor()] };
  client = new WaypointClient('https://localhost:1235', this.credentails, null);
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    api: ApiService;
  }
}
