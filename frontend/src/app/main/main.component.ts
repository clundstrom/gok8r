import {Component, OnInit} from '@angular/core';
import {MatSnackBar} from "@angular/material/snack-bar";
import {ApiService} from "../services/api.service";
import {AppConfigService} from "../services/app-config.service";
import {Observable} from "rxjs";
import {WebsockService} from "../services/websock.service";

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent implements OnInit {

  private SNACKBAR_UPTIME_MS = 3000;
  private readonly STATUS = 'Status: ';
  private readonly TASK_TRIGGER_FAILED_MSG = 'Could not connect to API';
  private readonly DISCONNECTED = 'Disconnected';
  private $CONNECTION: Observable<string> | undefined;
  public VERSION = "";
  spinner = false;

  constructor(private wss: WebsockService, private api: ApiService, private snackBar: MatSnackBar, private config: AppConfigService) {
  }

  ngOnInit(): void {
    // Initiate version
    this.VERSION = this.config.config.version;

    // Bind streams
    this.spinner = true;
    this.$CONNECTION = this.api.bindStream();
    this.$CONNECTION.subscribe(
      (msg: any) => this.snackBar.open(this.STATUS + msg.message)._dismissAfter(this.SNACKBAR_UPTIME_MS),
      (err: any) => {
        if (err.type == "close"){
          this.snackBar.open(this.DISCONNECTED)._dismissAfter(this.SNACKBAR_UPTIME_MS);
        }
        else{
          this.snackBar.open(this.TASK_TRIGGER_FAILED_MSG)._dismissAfter(this.SNACKBAR_UPTIME_MS);
          console.log(err);
        }
        this.spinner = false;
      });
  }

  triggerSSE() {
    this.api.sendSSE().subscribe((res) => res)
  }

  triggerWs(){
    this.api.sendWs().subscribe((res) => res)
  }

  queueJob() {
    this.api.queueJob(30).subscribe();
  }

  displayConfig() {
    this.snackBar.open("Current API Host: " + this.config.getApiHost())._dismissAfter(this.SNACKBAR_UPTIME_MS);
  }
}
