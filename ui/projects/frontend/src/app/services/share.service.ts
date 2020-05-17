import { Injectable } from '@angular/core';
import { PathService } from '../services/path.service';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

import { Share } from '../models/share';
import { Photo } from '../models/photo';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class ShareService {
  constructor(
    private readonly pathService: PathService,
    private http: HttpClient
  ) {}

  forSlug(slug: string): Observable<Share> {
    const url = this.pathService.shareBySlug(slug);

    return this.http.get<ShareAndPhotoResponse>(url).pipe(
      map((resp) => {
        const share = new Share();
        share.id = resp.share.id;
        share.slug = resp.share.slug;
        share.photos = resp.photos.map((photo) => ({
          ...photo,
          renditions: photo.renditions.map((r) => ({ ...r, url: this.pathService.renditionBySlug(slug, r.id) })),
        }));

        return share;
      })
    );
  }
}

export class ShareAndPhotoResponse {
  share: ShareResponse;
  photos: Array<Photo>;
}

export class ShareResponse {
  id: number;
  slug: string;
}
