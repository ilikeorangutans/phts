import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { distinctUntilChanged } from 'rxjs/operators';

import { AuthService } from './auth.service';

@Injectable({
  providedIn: 'root',
})
export class SessionService {
  private readonly _hasSession = new BehaviorSubject(false);

  readonly hasSession = this._hasSession
    .asObservable()
    .pipe(distinctUntilChanged());

  constructor(private readonly authService: AuthService) {}

  destroy() {
    this.authService.logout();
    this._hasSession.next(false);
  }

  start(username: string, password: string) {
    this.authService.authenticate(username, password).subscribe((status) => {
      if (status.authenticated) {
        this._hasSession.next(true);
      } else {
        this._hasSession.next(false);
      }
    });
  }
}
