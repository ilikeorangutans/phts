import { Collection } from './../models/collection';
import { Share } from './../models/share';
import { Photo } from './../models/photo';
import { Http } from '@angular/http';
import { PathService } from './path.service';
import { Injectable, Input } from '@angular/core';

@Injectable()
export class ShareService {

  constructor(
    private pathService: PathService,
    private http: Http
  ) { }

  listForPhoto(collection: Collection, photo: Photo): Promise<Array<Share>> {
    const url = this.pathService.photoShares(collection, photo.id);
    console.log(url);

    return this.http.get(url)
      .toPromise()
      .then(response => {
        const raw = response.json() as Array<Share>;

        return raw.map((share) => {
          share.createdAt = new Date(share.createdAt);
          share.updatedAt = new Date(share.updatedAt);

          return share;
        });
      });
  }

  save(collection: Collection, photo: Photo, share: Share): Promise<Share> {
    console.log(share);
    share.photoID = photo.id;

    const url = this.pathService.photoShares(collection, photo.id);
    return this.http.post(url, share)
      .toPromise()
      .then(response => {
        return response.json() as Share;
      });
  }
}
