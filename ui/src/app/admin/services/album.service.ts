import { HttpClient } from '@angular/common/http';
import { PathService } from './path.service';
import { Injectable } from '@angular/core';

import { Collection } from './../models/collection';
import { Album } from './../models/album';
import { Photo } from '../models/photo';

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

  addPhotos(collection: Collection, album: Album, photos: Array<Photo>) {
    console.log(`adding ${photos.length} photos to album ${album.name}`);
    const url = this.pathService.albumPhotos(collection, album);
    console.log(url);

    const photoIDs = photos.map(p => p.id);
    const submission = new PhotoSubmission(album.id, photoIDs);

    this.http.post(url, submission).toPromise().then(x => console.log("success"));
  }

  details(collection: Collection, album: Album) {
    const url = this.pathService.albumDetails(collection, album);
    console.log(url);
  }
}

interface PaginatedAlbums {
  data: Album[];
}

class PhotoSubmission {
  constructor(
    private albumID: number,
    private photoIDs: Array<number>
  ) {}
}
