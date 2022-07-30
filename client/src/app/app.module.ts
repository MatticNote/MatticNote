import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { NgProgressModule } from 'ngx-progressbar';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import {ProgressService} from "./service/progress.service";
import {FormsModule} from "@angular/forms";
import {NgIconsModule} from "@ng-icons/core";
import { BootstrapBell, BootstrapHouseDoor, BootstrapCompass } from '@ng-icons/bootstrap-icons';
import {CookieService} from "ngx-cookie-service";
import {MNAPIService} from "./service/mnapi.service";
import {HttpClientModule} from "@angular/common/http";
import {ToastContainerModule, ToastrModule} from "ngx-toastr";
import {ToastComponent} from "./component/toast/toast.component";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";

@NgModule({
  declarations: [
    AppComponent,
    ToastComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    NgProgressModule,
    FormsModule,
    HttpClientModule,
    NgIconsModule.withIcons({
      BootstrapBell,
      BootstrapHouseDoor,
      BootstrapCompass,
    }),
    BrowserAnimationsModule,
    ToastrModule.forRoot({
      toastComponent: ToastComponent,
      timeOut: 5000,
      extendedTimeOut: 10000,
      iconClasses: {
        info: 'alert-info',
        success: 'alert-success',
        warning: 'alert-warning',
        error: 'alert-error',
      },
      toastClass: 'alert',
      titleClass: '',
      messageClass: '',
      positionClass: 'toast',
    }),
    ToastContainerModule,
  ],
  providers: [
    ProgressService,
    CookieService,
    MNAPIService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
