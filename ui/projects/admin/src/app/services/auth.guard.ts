import { Injectable } from '@angular/core';
import {
  CanActivate,
  ActivatedRouteSnapshot,
  RouterStateSnapshot,
} from '@angular/router';
import { Observable } from 'rxjs';

import { Router } from '@angular/router';
import { CanActivateChild } from '@angular/router';
import { SessionService } from './session.service';

@Injectable({ providedIn: 'root' })
export class AuthGuard implements CanActivate, CanActivateChild {
  private isLoggedIn = false;

  constructor(private router: Router, private sessionService: SessionService) {
    this.sessionService.hasSession.subscribe(
      (loggedIn) => (this.isLoggedIn = loggedIn)
    );
  }

  canActivate(
    _next: ActivatedRouteSnapshot,
    _state: RouterStateSnapshot
  ): Observable<boolean> | Promise<boolean> | boolean {
    return this.checkAuth();
  }

  canActivateChild(
    _childRoute: ActivatedRouteSnapshot,
    _state: RouterStateSnapshot
  ): boolean | Observable<boolean> | Promise<boolean> {
    return this.checkAuth();
  }

  private checkAuth(): boolean {
    if (!this.isLoggedIn) {
      this.router.navigate(['login']);
    }

    return this.isLoggedIn;
  }
}
