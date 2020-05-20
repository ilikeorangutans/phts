import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { PathService } from './path.service';

import { map } from 'rxjs/operators';

class AuthResponse {
  sessionID: string;
  errors: Array<string>;
}

export class AuthStatus {
  static fromAuthResponse(resp: AuthResponse): AuthStatus {
    return new AuthStatus(resp.errors.length === 0, resp.sessionID);
  }
  constructor(readonly authenticated: boolean, readonly sessionID: string) {}
}

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  constructor(
    private readonly http: HttpClient,
    private readonly pathService: PathService
  ) {}

  authenticate(username: string, password: string): Observable<AuthStatus> {
    const url = this.pathService.sessionCreate();
    const headers = new HttpHeaders();
    headers.set('content-type', 'application/json');
    return this.http
      .post<AuthResponse>(
        url,
        { username, password },
        { withCredentials: true, headers }
      )
      .pipe(
        map((resp) => {
          return AuthStatus.fromAuthResponse(resp);
        })
      );
  }

  logout() {
    // TODO destroy session token
    const url = this.pathService.sessionDestroy();
    this.http.post(url, {}, { withCredentials: true }).subscribe();
  }
}
