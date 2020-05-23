import { PathService } from './path.service';
import { Injectable } from '@angular/core';
import {
  HttpEvent,
  HttpRequest,
  HttpInterceptor,
  HttpHandler,
} from '@angular/common/http';

import { SessionService } from './session.service';
import { Observable } from 'rxjs';

@Injectable()
export class JWTInterceptor implements HttpInterceptor {
  constructor(
    private sessionService: SessionService,
    private pathservice: PathService
  ) {}

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    if (req.url === this.pathservice.authenticate()) {
      return next.handle(req);
    }

    const authReq = req.clone({
      headers: req.headers.append('X-JWT', this.sessionService.getJWT()),
      withCredentials: true, // dirty hack so we always send cookies because the auth scheme here is a mess
    });

    return next.handle(authReq);
  }
}
