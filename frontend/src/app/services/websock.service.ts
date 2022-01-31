import {Injectable} from '@angular/core';
import {Observable} from "rxjs";
import {EventPublisher} from "../intefaces/eventpublisher";
import {webSocket,WebSocketSubject} from "rxjs/webSocket";
import {AppConfigService} from "./app-config.service";
import * as uuid from 'uuid';

@Injectable({
  providedIn: 'root'
})
export class WebsockService implements EventPublisher {

  constructor(private config: AppConfigService) {}

  private readonly socketId = uuid.v4();
  private $socket: WebSocketSubject<any> | undefined;

  public bind(uri: any): Observable<any> {
    this.$socket = webSocket(this.config.getWebSocket() + "?id=" + this.socketId);
    return this.$socket.asObservable();
  }

  getSocketId(){
    return this.socketId;
  }
  sendMessage(message: string) {
    console.log(message);
    if(this.$socket){
      this.$socket.next(message);
    }
  }
}
