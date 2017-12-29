import { Injectable } from '@angular/core';

import { AuthResponse } from '../auth.service';
import { User } from './../models/user';

@Injectable()
export class SessionService {

  private user: User;

  constructor() { }

  getJWT(): string {
    return localStorage.getItem('AuthService.jwt');
  }

  getUser(): User {
    return this.user;
  }

  login(auth: AuthResponse) {
    localStorage.setItem('AuthService.userID', auth.id.toString(10));
    localStorage.setItem('AuthService.userEmail', auth.email);
    localStorage.setItem('AuthService.jwt', auth.jwt);

    this.user = new User();
    this.user.id = auth.id;
    this.user.email = auth.email;
  }

  isLoggedIn(): Boolean {
    return this.getJWT() !== null;
  }

  logout() {
    localStorage.removeItem('AuthService.jwt');
  }
}

