import { Component } from '@angular/core';
import {Toast, ToastPackage, ToastrService} from "ngx-toastr";

@Component({
  selector: '[toast-component]',
  templateUrl: './toast.component.html',
  styleUrls: ['./toast.component.scss']
})
export class ToastComponent extends Toast {
  constructor(
    protected ts: ToastrService,
    public tp: ToastPackage,
  ) {
    super(ts, tp);
  }
}
