import {Injectable, NgZone} from '@angular/core';
import {HttpClient, HttpErrorResponse} from "@angular/common/http";
import {AppConfigService} from "./app-config.service";
import {EventPublisher} from "../intefaces/eventpublisher";
import {SseService} from "./sse.service";
import {catchError} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  publishers: Array<EventPublisher> = [];

  constructor(private http: HttpClient, private _zone: NgZone, private config: AppConfigService, private sse: SseService) {
  }

  private static handleError(error: HttpErrorResponse) {
    if (error.status === 0) {
      console.error('An error occurred:', error.error);
    } else {
      console.error(
        `Backend returned code ${error.status}, body was: `, error.error);
    }
    // Return an observable with a user-facing error message.
    return 'Something bad happened; please try again later.';
  }

  get(uri: string) {
    return this.http.get(this.config.getApiHost() + uri);
  }

  getHost() {
    return this.get("/api/v1/sendmessage").pipe(catchError(ApiService.handleError))
  }

  bindStream() {
    return this.sse.bind("/api/v1/stream");
  }
}
