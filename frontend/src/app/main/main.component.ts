import { Component, OnInit } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {MatSnackBar} from "@angular/material/snack-bar";

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent implements OnInit {

  constructor(private http: HttpClient, private snackBar: MatSnackBar) { }

  ngOnInit(): void {
  }

  spinner = false;

  btnTrigger() {
    this.snackBar.open('Task triggered.')._dismissAfter(3000);
    this.spinner = true;
  }
}
