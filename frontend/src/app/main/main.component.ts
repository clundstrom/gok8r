import {Component, OnInit} from '@angular/core';
import {MatSnackBar} from "@angular/material/snack-bar";
import {ApiService} from "../services/api.service";
import {AppConfigService} from "../services/app-config.service";

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent implements OnInit {

  private SNACKBAR_UPTIME_MS = 3000;
  private readonly TASK_TRIGGER_MSG = 'Connected to API: ';
  private readonly TASK_TRIGGER_FAILED_MSG = 'Could not connect to API';
  spinner = false;

  constructor(private api: ApiService, private snackBar: MatSnackBar, private config: AppConfigService) {
  }

  ngOnInit(): void {
  }

  btnTrigger() {
    this.spinner = true;

    this.api.getSSE("/api/v1/handshake").subscribe(
      res => {
        this.snackBar.open(this.TASK_TRIGGER_MSG + res)._dismissAfter(this.SNACKBAR_UPTIME_MS)
      },
      err => {
        this.snackBar.open(this.TASK_TRIGGER_FAILED_MSG)._dismissAfter(this.SNACKBAR_UPTIME_MS);
        console.log(err);
        this.spinner = false;
      },
    );
  }

  displayConfig() {
    this.snackBar.open("Current API Host: " + this.config.config.apiHostUrl)._dismissAfter(this.SNACKBAR_UPTIME_MS);
  }
}
