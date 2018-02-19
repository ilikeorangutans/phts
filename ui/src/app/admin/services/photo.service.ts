import { Observable } from 'rxjs/Observable';
import { Album } from './../models/album';
import { RenditionConfiguration } from './../models/rendition-configuration';
import { Paginator } from './../models/paginator';
import { HttpClient } from '@angular/common/http';
import { PathService } from './path.service';
import { Injectable } from '@angular/core';
import { Photo } from '../models/photo';
import { Collection } from '../models/collection';

@Injectable()
export class PhotoService {

  constructor(
    private pathService: PathService,
    private http: HttpClient
  ) { }

  list(collection: Collection, paginator: Paginator): Promise<Array<Photo>> {
    const path = this.pathService.listPhotos(collection);

    const url = `${path}?${paginator.toQueryString()}`;

    return this.http
      .get<PaginatedPhotos>(url)
      .toPromise()
      .then((response) => {
        return response.data.map((photo) => {
          photo.collection = collection;
          photo.updatedAt = new Date(photo.updatedAt);
          photo.createdAt = new Date(photo.createdAt);
          if (photo.takenAt) {
            photo.takenAt = new Date(photo.takenAt);
          }

          return photo;
        });
      });
  }

  listAlbum(collection: Collection, album: Album, paginator: Paginator): Observable<Array<Photo>> {
    const path = this.pathService.albumPhotos(collection, album);

    const url = `${path}?${paginator.toQueryString()}`;

    console.log(url);

    return this.http
      .get<PaginatedPhotos>(url)
      .map((response) => {
        return response.data.map((photo) => {
          photo.collection = collection;
          photo.updatedAt = new Date(photo.updatedAt);
          photo.createdAt = new Date(photo.createdAt);
          if (photo.takenAt) {
            photo.takenAt = new Date(photo.takenAt);
          }

          return photo;
        });
      });
  }

  byID(collection: Collection, photoID: number, renditionConfigurations: RenditionConfiguration[]): Promise<Photo> {
    let queryString = '';
    if (renditionConfigurations.length > 0) {
      queryString = `?rendition-configuration-ids=${renditionConfigurations.map((c => c.id)).join(',')}`;
    }
    const url = `${this.pathService.showPhoto(collection, photoID)}${queryString}`;
    return this.http
      .get<Photo>(url)
      .toPromise();
  }

  recentPhotos(collection: Collection, renditionConfigurations: RenditionConfiguration[]): Promise<Photo[]> {
    let queryString = '';
    if (renditionConfigurations.length > 0) {
      queryString = `?rendition-configuration-ids=${renditionConfigurations.map((c => c.id)).join(',')}`;
    }
    const url = `${this.pathService.recentPhotos(collection)}${queryString}`;
    return this.http
      .get<PaginatedPhotos>(url)
      .toPromise()
      .then((response) => {
        return response.data.map((photo) => {
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
      .catch((e) => {
        return Promise.reject(e);
      });
  }

  upload(collection: Collection, file: File): Promise<Photo> {
    const url = this.pathService.uploadPhoto(collection);
    const formdata = new FormData();
    formdata.append('image', file, file.name);
    return this.http.post<Photo>(url, formdata)
      .toPromise();
  }

  publish(collection: Collection, photo: Photo): Promise<Photo> {
    return Promise.reject('Implement me!');
  }
}

interface PaginatedPhotos {
  data: Photo[];
}
