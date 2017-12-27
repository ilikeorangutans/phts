import { Injectable } from '@angular/core';

@Injectable()
export class AuthService {

  private loggedIn: Boolean = false;

  constructor() { }

  login() {
    // TODO this is obviously work in progress
    console.log('I am not really logging you in. This is fake. There is no login.');
    this.loggedIn = true;
  }

  logout() {
    this.loggedIn = false;
  }

  isLoggedIn(): Boolean {
    return this.loggedIn;
  }

}
