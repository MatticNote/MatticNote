import {Component, OnInit} from '@angular/core';
import {ProgressService} from "./service/progress.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  isDrawerOpen: boolean = false;

  constructor(private ps: ProgressService) {
  }

  ngOnInit() {
    this.ps.init();
  }
}
