import {Injectable} from '@angular/core';
import {Observable} from "rxjs";
import {EventPublisher} from "../intefaces/eventpublisher";

@Injectable({
  providedIn: 'root'
})
export class WebsockService implements EventPublisher {

  constructor() {
  }

  bind(uri: any): Observable<any> {
    return new Observable<any>((observer) => {
    });
  }
}
