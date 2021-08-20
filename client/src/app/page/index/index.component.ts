import { Component, OnInit } from '@angular/core';
import {URLTrusterService} from "../../service/url-truster/url-truster.service";

@Component({
  selector: 'app-index',
  templateUrl: './index.component.html',
  styleUrls: ['./index.component.scss']
})
export class IndexComponent implements OnInit {

  constructor(
    public ut: URLTrusterService,
  ) { }

  ngOnInit(): void {
  }

}
