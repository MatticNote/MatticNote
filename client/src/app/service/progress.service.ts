import { Injectable } from '@angular/core';
import {NgProgress, NgProgressRef} from "ngx-progressbar";

@Injectable({
  providedIn: 'root'
})
export class ProgressService {
  private progressRef?: NgProgressRef;

  constructor(private progress: NgProgress) { }

  init() {
    this.progressRef = this.progress.ref('progress');
  }

  start() {
    this.progressRef?.start();
  }

  complete() {
    this.progressRef?.complete();
  }
}
