import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { distinctUntilChanged, tap, map } from 'rxjs/operators';

import { CookieService } from 'ngx-cookie-service';

import { AuthService } from './auth.service';

export class SessionStatus {
  constructor(readonly message: string) {}
}

@Injectable({
  providedIn: 'root',
})
export class SessionService {
  readonly SESSION_COOKIE_NAME = 'PHTS_SERVICES_INTERNAL_SESSION_ID';

  private readonly _hasSession = new BehaviorSubject(false);

  readonly hasSession = this._hasSession
    .asObservable()
    .pipe(distinctUntilChanged());

  constructor(
    private readonly authService: AuthService,
    private readonly cookies: CookieService
  ) {
    if (this.cookies.check(this.SESSION_COOKIE_NAME)) {
      this._hasSession.next(true);
    }
  }

  destroy() {
    this.authService.logout();
    this._hasSession.next(false);
  }

  start(username: string, password: string): Observable<SessionStatus> {
    return this.authService.authenticate(username, password).pipe(
      tap((status) => {
        if (status.authenticated) {
          this._hasSession.next(true);
        } else {
          this._hasSession.next(false);
        }
      }),
      map((_) => new SessionStatus('TODO insert message'))
    );
  }
}
