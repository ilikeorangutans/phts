import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
} from '@angular/common/http';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { SessionService } from './session.service';

/**
 * Checks incoming responses from the server.
 * TODO check responses if they are auth failures. If so, end the session.
 */
@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  constructor(private readonly sessionService: SessionService) {}

  intercept(
    request: HttpRequest<unknown>,
    next: HttpHandler
  ): Observable<HttpEvent<unknown>> {
    return next.handle(request).pipe(
      tap(
        (_) => {},
        (error) => {
          // TODO would be nice if we'd notify the user why his session has been ended
          if (error.status === 403 || error.status === 401) {
            this.sessionService.destroy();
          }
        }
      )
    );
  }
}
