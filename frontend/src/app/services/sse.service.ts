import {Injectable, NgZone} from '@angular/core';
import {Observable} from "rxjs";
import {EventPublisher} from "../intefaces/eventpublisher";

@Injectable({
  providedIn: 'root'
})
export class SseService implements EventPublisher {

  constructor(private _zone: NgZone) {
  }

  private readonly CONNECTED = "Connected";

  bind(uri: any): Observable<any> {
    return new Observable<any>((observer) => {
      let eventSource = new EventSource(uri);

      eventSource.onopen = (event) => {
        console.log(this.CONNECTED);
      };

      eventSource.onmessage = (event) => {
        this._zone.run(() => observer.next(event.data));
      };

      eventSource.onerror = (error) => {
        if (eventSource.readyState === 0) {
          console.log('The stream has been closed by the server.');
          observer.complete();
        } else {
          observer.error(error);
        }
      }
    });
  }
}
