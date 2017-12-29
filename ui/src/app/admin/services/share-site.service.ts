import { Observable } from 'rxjs/Observable';
import { AuthService } from './../auth.service';
import { PathService } from './path.service';
import { Injectable } from '@angular/core';
import { ShareSite } from '../models/share-site';
import { HttpClient, HttpEvent, HttpRequest, HttpInterceptor, HttpHandler } from '@angular/common/http';
import { SessionService } from './session.service';

@Injectable()
export class ShareSiteService {

  constructor(
    private pathService: PathService,
    private http: HttpClient
  ) { }

  list(): Promise<Array<ShareSite>> {
    const url = this.pathService.shareSites();

    return this.http
      .get<Array<ShareSite>>(url)
      .toPromise()
      .then((response) => {
        const records = response;

        return records
          .map(r => {
            r.createdAt = new Date(r.createdAt);
            r.updatedAt = new Date(r.updatedAt);
            return r;
          });
      });
  }

  save(shareSite: ShareSite): Promise<ShareSite> {
    const url = this.pathService.shareSites();
    return this.http.post<ShareSite>(url, shareSite)
      .toPromise()
      .then((response) => {
        return response;
      });
  }
}

@Injectable()
export class JWTInterceptor implements HttpInterceptor {

  constructor(
    private sessionService: SessionService
  ) {}

  intercept(req: HttpRequest<any>, next: HttpHandler):
    Observable<HttpEvent<any>> {

    console.log('interceptor blargh', req.url);
    // TODO this is pretty shitty because it will catch all requests
    if (req.url.endsWith('/admin/api/authenticate')) {
      return next.handle(req);
    }

    console.log('interceptor adding token', this.sessionService.getJWT());

    const authReq = req.clone({headers: req.headers.append('X-JWT', this.sessionService.getJWT())});

    return next.handle(authReq);
  }
}
