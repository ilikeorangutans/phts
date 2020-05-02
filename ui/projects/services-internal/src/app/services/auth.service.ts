import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly _authenticated = new BehaviorSubject(false);

  /**
   * Observable of the authenticated state.
   */
  readonly authenticated = this._authenticated.asObservable();

  constructor() {}

  authenticate(username: string, password: string) {
    console.log(`logging in with ${username} and ${password}`)
    this._authenticated.next(true);
  }

  logout() {
    this._authenticated.next(false);
    // TODO destroy session token
  }
}
