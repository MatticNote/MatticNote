import { Injectable } from '@angular/core';
import { environment } from "../../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class URLTrusterService {

  constructor() { }

  openExternalURL(event: MouseEvent): void {
    event.preventDefault();
    const elem = event.target as HTMLLinkElement;
    if (!this.isTrusted(elem.href) && !window.confirm(`${elem.href} is not trusted domain. Are you sure to access?`)) {
      return;
    }
    open(elem.href, '_blank');
  }

  isTrusted(href: string): boolean {
    try {
      const targetUrl = new URL(href);
      const res = environment.defaultTrustedDomain.find(
        val => targetUrl.host.match(new RegExp(`^(.+\.)*${val.replace(/\./g, '\\.')}$`, 'i'))
      );
      if (res !== undefined) {
        return true;
      }
      return false;
    } catch (e) {
      if (e instanceof TypeError) {
        return false;
      }
      console.error(e);
      return false;
    }
  }

}
