import {Injectable, NgZone} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  constructor(private http: HttpClient, private _zone: NgZone) { }

  getTaskResult(){
    return this.http.get('/api/v1/', {responseType: 'text'})
  }

  getSSE(uri: any): Observable<string> {
    return new Observable<string>((observer) => {
      let eventSource = new EventSource(uri);
      eventSource.onmessage = (event) => {
        observer.next(event.data);
      };
      eventSource.onerror = (error) => {
        if(eventSource.readyState === 0) {
          console.log('The stream has been closed by the server.');
          eventSource.close();
          observer.complete();
        } else {
          observer.error('EventSource error: ' + error);
        }
      }
    });
  }
}
