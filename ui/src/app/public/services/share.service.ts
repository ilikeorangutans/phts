import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { PathService } from './path.service';
import { Share } from './../models/share';
import { Photo } from './../models/photo';

@Injectable()
export class ShareService {

  constructor(
    private http: HttpClient,
    private pathService: PathService
  ) { }

  forSlug(slug: string): Promise<Share> {
    const url = this.pathService.shareBySlug(slug);
    return this.http
      .get<ShareAndPhotoResponse>(url)
      .toPromise()
      .then(resp => {
        const share = new Share();
        share.id = resp.share.id;
        share.slug = resp.share.slug;
        share.photos = resp.photos.map(p => {
          const photo = new Photo();
          photo.id = p.id;
          photo.renditions = p.renditions;
          return photo;
        });

        return share;
      });
  }
}

class ShareAndPhotoResponse {
  share: ShareResponse;
  photos: Array<Photo>;
}

class ShareResponse {
  id: number;
  slug: string;
}

