import { Injectable } from '@angular/core';
import {CookieService} from "ngx-cookie-service";
import {User} from "../interface/user";
import {HttpClient, HttpErrorResponse} from "@angular/common/http";
import {Observable} from "rxjs";
import {ProgressService} from "./progress.service";
import {APIError} from "../interface/api-error";
import {ToastrService} from "ngx-toastr";

@Injectable({
  providedIn: 'root'
})
export class MNAPIService {
  private token?: string;
  private _user?: User;

  constructor(
    private cs: CookieService,
    private hs: HttpClient,
    private ps: ProgressService,
    private ts: ToastrService,
  ) { }

  get user(): User | undefined {
    return this._user;
  }

  init() {
    this.token = this.cs.get('mn_token');
    this.ps.start();
    this.usersMe()
      .subscribe({
        next: value => {
          this._user = value;
        },
        error: (err: HttpErrorResponse) => {
          if (err.status === 0) {
            this.ts.error('Cannot reach to server.', 'Network Error')
          } else {
            const api_error = err.error as APIError;
            switch (api_error.code) {
              case 'UNAUTHORIZED':
                location.href = '/account/login';
                break;
              default:
                console.error(err);
            }
          }
          this.ps.complete();
        },
        complete: () => {
          this.ps.complete();
        }
      });
  }

  usersMe(): Observable<User> {
    return this.hs.get('/api/v1/users/me', {
      headers: {
        'Authorization': `token ${this.token}`
      },
      withCredentials: false,
    }) as Observable<User>;
  }
}
