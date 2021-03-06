import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { SessionService } from './session.service';
import { PathService } from './path.service';

export class Credentials {
  username: string;
  password: string;
}

export class AuthResponse {
  email: string;
  id: number;
  jwt: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  constructor(
    private http: HttpClient,
    private sessionService: SessionService,
    private pathService: PathService
  ) {}

  authenticate(credentials: Credentials): Promise<boolean> {
    const url = this.pathService.authenticate();
    return this.http
      .post<AuthResponse>(url, credentials, { withCredentials: true })
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
