import { HttpClient } from '@angular/common/http';
import { PathService } from './path.service';
import { Injectable } from '@angular/core';

import { Collection } from './../models/collection';
import { Album } from './../models/album';

@Injectable()
export class AlbumService {

  constructor(
    private pathService: PathService,
    private http: HttpClient
  ) { }

  list(collection: Collection): Promise<Array<Album>> {
    const url = this.pathService.albumBase(collection);
    return this.http.get<PaginatedAlbums>(url)
      .toPromise()
      .then(resp => {
        return resp.data.map(album => {
          album.createdAt = new Date(album.createdAt);
          album.updatedAt = new Date(album.updatedAt);
          return album;
        });
      });
  }

  save(collection: Collection, album: Album): Promise<Album> {
    const url = this.pathService.albumBase(collection);
    return this.http.post<Album>(url, album).toPromise();
  }
}

interface PaginatedAlbums {
  data: Album[];
}
