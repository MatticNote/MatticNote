import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { NgProgressModule } from 'ngx-progressbar';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import {ProgressService} from "./service/progress.service";
import {FormsModule} from "@angular/forms";
import {NgIconsModule} from "@ng-icons/core";
import { BootstrapBell, BootstrapHouseDoor, BootstrapCompass } from '@ng-icons/bootstrap-icons';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    NgProgressModule,
    FormsModule,
    NgIconsModule.withIcons({
      BootstrapBell,
      BootstrapHouseDoor,
      BootstrapCompass,
    }),
  ],
  providers: [ProgressService],
  bootstrap: [AppComponent]
})
export class AppModule { }
