import {Injectable, Injector} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class AppConfigService {

  private appConfig: any;
  private localApiHost = "http://localhost:4200";

  constructor (private injector: Injector) { }

  loadAppConfig() {
    let http = this.injector.get(HttpClient);
    return http.get('/assets/app-config.json')
      .toPromise()
      .then(data => {
        this.appConfig = data;
      })
  }

  get config() {
    return this.appConfig;
  }

  getApiHost(){
    if(!environment.production){
      return this.localApiHost;
    }
    return this.appConfig.apiHostUrl;
  }
}
