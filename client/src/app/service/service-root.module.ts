import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {SessionService} from "./session/session.service";
import {URLTrusterService} from "./url-truster/url-truster.service";



@NgModule({
  declarations: [],
  imports: [
    CommonModule
  ],
  providers: [
    SessionService,
    URLTrusterService,
  ]
})
export class ServiceRootModule { }
