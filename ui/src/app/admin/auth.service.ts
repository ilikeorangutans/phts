import { Injectable } from '@angular/core';
import 'rxjs/add/operator/map';
import { HttpClient } from '@angular/common/http';
import { SessionService } from './services/session.service';

export class Credentials {
  username: string;
  password: string;
}

export class AuthResponse {
  email: string;
  id: number;
  jwt: string;
}

@Injectable()
export class AuthService {

  constructor(
    private http: HttpClient,
    private sessionService: SessionService
  ) { }

  authenticate(credentials: Credentials): Promise<Boolean> {
    return this.http.post<AuthResponse>('http://localhost:8080/admin/api/authenticate', credentials)
      .toPromise()
      .then((resp) => {
        this.sessionService.login(resp);
        return true;
      })
      .catch((e) => {
        console.log('login failed', e);
        return false;
      });
  }

  // TODO need a logout method here to destroy session on server

}
