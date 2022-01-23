import {Injectable, NgZone} from '@angular/core';
import {HttpClient, HttpErrorResponse} from "@angular/common/http";
import {AppConfigService} from "./app-config.service";
import {SseService} from "./sse.service";
import {catchError} from 'rxjs/operators';
import {WebsockService} from "./websock.service";

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  constructor(private http: HttpClient, private _zone: NgZone,
              private config: AppConfigService,
              private ws: WebsockService,
              private sse: SseService) {
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

  sendMessage() {
    return this.get("/api/v1/sendmessage").pipe(catchError(ApiService.handleError))
  }

  bindStream() {
    let $obs;
    try {
      $obs = this.ws.bind(this.config.getApiHost() + "/api/v1/socket");
    } catch (e) {
      console.warn("Could not bind to websocket, falling back on SSE.");
      $obs = this.sse.bind(this.config.getApiHost() + "/api/v1/stream");
    }
    return $obs;
  }
}
