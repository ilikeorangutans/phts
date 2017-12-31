import { Share } from './../models/share';
import { Photo } from './../models/photo';
import { Collection } from './../models/collection';
import { HttpClient } from '@angular/common/http';
import { PathService } from './path.service';
import { Injectable } from '@angular/core';

@Injectable()
export class ShareService {

  constructor(
    private pathService: PathService,
    private http: HttpClient
  ) { }

  listForPhoto(collection: Collection, photo: Photo): Promise<Array<Share>> {
    const url = this.pathService.photoShares(collection, photo.id);

    return this.http.get<Array<Share>>(url)
      .toPromise()
      .then(response => {
        return response.map((share) => {
          share.createdAt = new Date(share.createdAt);
          share.updatedAt = new Date(share.updatedAt);

          return share;
        });
      });
  }

  save(collection: Collection, photo: Photo, share: Share): Promise<Share> {
    share.photoID = photo.id;

    const url = this.pathService.photoShares(collection, photo.id);
    return this.http.post<Share>(url, share)
      .toPromise();
  }
}
