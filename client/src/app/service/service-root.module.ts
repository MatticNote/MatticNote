import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {SessionService} from "./session/session.service";



@NgModule({
  declarations: [],
  imports: [
    CommonModule
  ],
  providers: [
    SessionService,
  ]
})
export class ServiceRootModule { }
