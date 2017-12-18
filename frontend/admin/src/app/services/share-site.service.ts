import { ShareSite } from './../models/share-site';
import { Http } from '@angular/http';
import { PathService } from './path.service';
import { Injectable } from '@angular/core';

@Injectable()
export class ShareSiteService {

  constructor(
    private pathService: PathService,
    private http: Http
  ) { }

  list(): Promise<Array<ShareSite>> {
    const url = this.pathService.shareSites();
    console.log(url);

    return this.http
      .get(url)
      .toPromise()
      .then((response) => {
        const records = response.json() as ShareSite[];

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
    return this.http.post(url, shareSite)
      .toPromise()
      .then((response) => {
        console.log(response);
        return response.json() as ShareSite;
      });
  }
}
