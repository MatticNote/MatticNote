import { Component, OnInit } from '@angular/core';
import {ProgressService} from "../../service/progress.service";

@Component({
  selector: 'app-index',
  templateUrl: './index.component.html',
  styleUrls: ['./index.component.scss']
})
export class IndexComponent implements OnInit {

  constructor(public ps: ProgressService) { }

  ngOnInit(): void {
  }

}
