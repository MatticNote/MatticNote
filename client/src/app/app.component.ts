import {Component, OnInit, ViewChild} from '@angular/core';
import {ProgressService} from "./service/progress.service";
import {MNAPIService} from "./service/mnapi.service";
import {ToastContainerDirective, ToastrService} from "ngx-toastr";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  @ViewChild(ToastContainerDirective, { static: true }) toastContainer?: ToastContainerDirective;

  isDrawerOpen: boolean = false;

  constructor(
    private ps: ProgressService,
    public api: MNAPIService,
    private ts: ToastrService,
  ) { }

  ngOnInit() {
    this.ts.overlayContainer = this.toastContainer;
    this.ps.init();
    this.api.init();
  }
}
