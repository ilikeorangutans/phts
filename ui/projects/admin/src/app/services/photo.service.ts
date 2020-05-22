import { Album } from './../models/album';
import { RenditionConfiguration } from './../models/rendition-configuration';
import { Paginator } from './../models/paginator';
import { HttpClient } from '@angular/common/http';
import { PathService } from './path.service';
import { Injectable } from '@angular/core';
import { Photo } from '../models/photo';
import { Collection } from '../models/collection';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable()
export class PhotoService {
  constructor(private pathService: PathService, private http: HttpClient) {}

  list(collection: Collection, paginator: Paginator): Observable<Array<Photo>> {
    const path = this.pathService.listPhotos(collection);
    const url = `${path}?${paginator.toQueryString()}`;

    return this.http.get<PaginatedPhotos>(url).pipe(
      map((response) => response.data),
      map((photos) => {
        return photos.map((photo) => {
          photo.collection = collection;
          photo.updatedAt = new Date(photo.updatedAt);
          photo.createdAt = new Date(photo.createdAt);
          if (photo.takenAt) {
            photo.takenAt = new Date(photo.takenAt);
          }
          return photo;
        });
      })
    );
  }

  listAlbum(
    collection: Collection,
    album: Album,
    paginator: Paginator
  ): Observable<Array<Photo>> {
    const path = this.pathService.albumPhotos(collection, album);

    const url = `${path}?${paginator.toQueryString()}`;

    return this.http.get<PaginatedPhotos>(url).pipe(
      map((response) => {
        return response.data.map((photo) => {
          photo.collection = collection;
          photo.updatedAt = new Date(photo.updatedAt);
          photo.createdAt = new Date(photo.createdAt);
          if (photo.takenAt) {
            photo.takenAt = new Date(photo.takenAt);
          }

          return photo;
        });
      })
    );
  }

  byID(
    collection: Collection,
    photoID: number,
    renditionConfigurations: RenditionConfiguration[]
  ): Observable<Photo> {
    let queryString = '';
    if (renditionConfigurations.length > 0) {
      queryString = `?rendition-configuration-ids=${renditionConfigurations
        .map((c) => c.id)
        .join(',')}`;
    }
    const url = `${this.pathService.showPhoto(
      collection,
      photoID
    )}${queryString}`;
    return this.http.get<Photo>(url);
  }

  recentPhotos(
    collection: Collection,
    renditionConfigurations: RenditionConfiguration[]
  ): Observable<Array<Photo>> {
    let queryString = '';
    if (renditionConfigurations.length > 0) {
      queryString = `?rendition-configuration-ids=${renditionConfigurations
        .map((c) => c.id)
        .join(',')}`;
    }
    const url = `${this.pathService.recentPhotos(collection)}${queryString}`;
    return this.http.get<PaginatedPhotos>(url).pipe(
      map((response) => response.data),
      map((photos) => {
        return photos.map((photo) => {
          photo.collection = collection;
          photo.updatedAt = new Date(photo.updatedAt);
          photo.createdAt = new Date(photo.createdAt);
          if (photo.takenAt) {
            photo.takenAt = new Date(photo.takenAt);
          }
          photo.renditions = photo.renditions.map((rendition) => {
            rendition.createdAt = new Date(rendition.createdAt);
            rendition.updatedAt = new Date(rendition.updatedAt);

            return rendition;
          });
          return photo;
        });
      })
    );
  }

  upload(collection: Collection, file: File): Observable<Photo> {
    const url = this.pathService.uploadPhoto(collection);
    const formdata = new FormData();
    formdata.append('image', file, file.name);
    return this.http.post<Photo>(url, formdata);
  }

  publish(collection: Collection, photo: Photo): Promise<Photo> {
    console.log(collection, photo);
    return Promise.reject('Implement me!');
  }

  delete(collection: Collection, photo: Photo) {
    const url = this.pathService.showPhoto(collection, photo.id);
    this.http.delete(url).subscribe((response) => console.log(response));
  }
}

interface PaginatedPhotos {
  data: Photo[];
}
