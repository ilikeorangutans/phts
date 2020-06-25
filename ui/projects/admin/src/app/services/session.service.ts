import { Injectable } from '@angular/core';

import { User } from './../models/user';
import { AuthResponse } from './auth.service';
import { BehaviorSubject } from 'rxjs';
import { distinctUntilChanged, map } from 'rxjs/operators';

export class ActiveSession {
  readonly hasSession: boolean = true;
  constructor(
    readonly user: User,
    readonly jwt: string // deprecating this
  ) {}
}

export class NoSession {
  readonly hasSession: boolean = false;
}

export type SessionStatus = ActiveSession | NoSession;

@Injectable({ providedIn: 'root' })
export class SessionService {
  private readonly status = new BehaviorSubject(new NoSession());

  /**
   * An observable that holds true if there's a session, or false if not.
   */
  readonly hasSession = this.status
    .asObservable()
    .pipe(map((status) => status.hasSession, distinctUntilChanged()));

  constructor() {
    this.status.next(this.statusFromLocalStorage());
  }

  statusFromLocalStorage(): SessionStatus {
    const hasJWT = this.getJWT() !== '';
    if (hasJWT) {
      return new ActiveSession(this.getUser(), this.getJWT());
    } else {
      return new NoSession();
    }
  }

  private user: User | null;

  getJWT(): string {
    const storedToken = localStorage.getItem('AuthService.jwt');
    if (storedToken === null) {
      return '';
    }
    return storedToken;
  }

  getUser(): User {
    if (this.user === undefined) {
      this.user = this.loadUser();
    }
    if (this.user === null) {
      throw 'no user';
    }
    return this.user;
  }

  loadUser(): User | null {
    const user = new User();
    const email = localStorage.getItem('AuthService.userEmail');
    if (email === null) {
      return null;
    }
    user.email = email;
    const userID = localStorage.getItem('AuthService.userID');
    if (userID === null) {
      return null;
    }
    user.id = +userID;
    return user;
  }

  login(auth: AuthResponse) {
    localStorage.setItem('AuthService.userID', auth.id.toString(10));
    localStorage.setItem('AuthService.userEmail', auth.email);
    localStorage.setItem('AuthService.jwt', auth.jwt);

    this.user = new User();
    this.user.id = auth.id;
    this.user.email = auth.email;
  }

  isLoggedIn(): boolean {
    return this.getJWT() !== null;
  }

  logout() {
    localStorage.removeItem('AuthService.jwt');
  }
}
