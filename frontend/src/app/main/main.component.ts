import {Component, OnInit} from '@angular/core';
import {MatSnackBar} from "@angular/material/snack-bar";
import {ApiService} from "../services/api.service";

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent implements OnInit {

  constructor(private api: ApiService, private snackBar: MatSnackBar) {
  }

  ngOnInit(): void {
  }

  spinner = false;

  btnTrigger() {
    this.spinner = true;
    this.api.getTaskResult().subscribe((res) => {
      if (res) {
        this.snackBar.open('Task triggered on host: ' + res)._dismissAfter(3000)
      }
    });

  }
}
