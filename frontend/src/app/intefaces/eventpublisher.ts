import {Observable} from "rxjs";

export interface EventPublisher {
  bind(uri: any): Observable<any>
}
