import {Injectable, NgZone} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private readonly CONNECTED = "Connected";

  constructor(private http: HttpClient, private _zone: NgZone) { }

  get(uri: string){
    return this.http.get(uri);
  }

  getSSE(uri: any): Observable<string> {
    return new Observable<string>((observer) => {
      let eventSource = new EventSource(uri);

      eventSource.onopen = (event) => {
        console.log(this.CONNECTED);
      };

      eventSource.onmessage = (event) => {
        this._zone.run(() => observer.next(event.data));
      };
      eventSource.onerror = (error) => {
        if(eventSource.readyState === 0) {
          console.log('The stream has been closed by the server.');
          observer.complete();
        } else {
          observer.error(error);
        }
      }
    });
  }
}
