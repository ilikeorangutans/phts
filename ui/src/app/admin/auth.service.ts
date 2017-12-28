import { User } from './models/user';
import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import 'rxjs/add/operator/map';

export class Credentials {
  username: string;
  password: string;
}

class AuthResponse {
  email: string;
  id: number;
  jwt: string;
}

@Injectable()
export class AuthService {

  private user: User;

  constructor(
    private http: Http
  ) { }

  authenticate(credentials: Credentials): Promise<Boolean> {
    return this.http.post('http://localhost:8080/admin/api/authenticate', credentials)
      .toPromise()
      .then((resp) => {
        if (resp.ok) {
          const authResp = resp.json() as AuthResponse;
          this.login(authResp);
          return true;
        }

        return false;
      })
      .catch((e) => {
        console.log('login failed', e);
        return false;
      });
  }

  login(auth: AuthResponse) {
    localStorage.setItem('AuthService.userID', auth.id.toString(10));
    localStorage.setItem('AuthService.userEmail', auth.email);
    localStorage.setItem('AuthService.jwt', auth.jwt);

    this.user = new User();
    this.user.id = auth.id;
    this.user.email = auth.email;
  }

  logout() {
    localStorage.removeItem('AuthService.jwt');
  }

  getJWT(): string {
    return localStorage.getItem('AuthService.jwt');
  }

  getUser(): User {
    return this.user;
  }

  isLoggedIn(): Boolean {
    return this.getJWT() !== null;
  }

}
