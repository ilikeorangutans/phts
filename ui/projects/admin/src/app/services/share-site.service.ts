import { PathService } from './path.service';
import { Injectable } from '@angular/core';
import { ShareSite } from '../models/share-site';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable()
export class ShareSiteService {
  constructor(private pathService: PathService, private http: HttpClient) {}

  list(): Observable<Array<ShareSite>> {
    const url = this.pathService.shareSites();

    return this.http.get<Array<ShareSite>>(url).pipe(
      map((records) => {
        return records.map((r) => {
          r.createdAt = new Date(r.createdAt);
          r.updatedAt = new Date(r.updatedAt);
          return r;
        });
      })
    );
  }

  save(shareSite: ShareSite): Observable<ShareSite> {
    const url = this.pathService.shareSites();
    return this.http.post<ShareSite>(url, shareSite);
  }
}
