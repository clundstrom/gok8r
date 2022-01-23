import {Injectable} from '@angular/core';
import {Observable} from "rxjs";
import {EventPublisher} from "../intefaces/eventpublisher";
import {webSocket,WebSocketSubject} from "rxjs/webSocket";
import {AppConfigService} from "./app-config.service";

@Injectable({
  providedIn: 'root'
})
export class WebsockService implements EventPublisher {

  constructor(private config: AppConfigService) {}

  private $socket: WebSocketSubject<any> | undefined;

  public bind(uri: any): Observable<any> {
    this.$socket = webSocket(this.config.getWebSocket());
    return this.$socket.asObservable();
  }

  sendMessage(message: string) {
    console.log(message);
    if(this.$socket){
      this.$socket.next(message);
    }
  }
}
