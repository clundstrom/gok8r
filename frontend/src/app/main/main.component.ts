import {Component, OnInit} from '@angular/core';
import {MatSnackBar} from "@angular/material/snack-bar";
import {ApiService} from "../services/api.service";

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent implements OnInit {

  private SNACKBAR_UPTIME_MS = 3000;
  private TASK_TRIGGER_MSG = 'Task triggered on host: ';
  spinner = false;
  eventLog: any = [];

  constructor(private api: ApiService, private snackBar: MatSnackBar) {
  }

  ngOnInit(): void {
  }

  btnTrigger() {
    this.spinner = true;

    this.api.getSSE("/api/v1/handshake").subscribe(res => {
      if (res) {
        this.snackBar.open(this.TASK_TRIGGER_MSG + res)._dismissAfter(this.SNACKBAR_UPTIME_MS);
        this.eventLog.push(res);
      }
    });
  }
}
