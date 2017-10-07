import { Injectable } from '@angular/core';
import { Http } from '@angular/http';

import { PathService } from './path.service';
import { Collection } from '../models/collection';
import { RenditionConfiguration } from '../models/rendition-configuration';
import { Photo } from '../models/photo';

@Injectable()
export class PhotoService {

  constructor(
    private pathService: PathService,
    private http: Http
  ) { }

  byID(collection: Collection, photoID: number, renditionConfigurations: RenditionConfiguration[]): Promise<Photo> {
    let queryString = '';
    if (renditionConfigurations.length > 0) {
      queryString = `?rendition-configuration-ids=${renditionConfigurations.map((c => c.id)).join(',')}`;
    }
    let url = `${this.pathService.showPhoto(collection, photoID)}${queryString}`;
    return this.http
      .get(url)
      .toPromise()
      .then((response) => {
        return response.json() as Photo;
      });
  }

  recentPhotos(collection: Collection, renditionConfigurations: RenditionConfiguration[]): Promise<Photo[]> {
    let queryString = '';
    if (renditionConfigurations.length > 0) {
      queryString = `?rendition-configuration-ids=${renditionConfigurations.map((c => c.id)).join(',')}`;
    }
    const url = `${this.pathService.recentPhotos(collection)}${queryString}`;
    return this.http
      .get(url)
      .toPromise()
      .then((response) => {
        const r = response.json() as PaginatedPhotos;

        return r.data.map((photo) => {
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

  upload(collection: Collection, file: File) {
    console.log("upload", collection, file);

    let url = this.pathService.uploadPhoto(collection);
    console.log("uploading photo to ", url)
    let formdata = new FormData();
    formdata.append("image", file, file.name);
    this.http.post(url, formdata)
      .toPromise()
      .then((response) => console.log(response))
      .catch((e) => console.log(e));
  }
}

interface PaginatedPhotos {
  data: Photo[];
}
