import { Injectable } from '@angular/core';
import { Http } from '@angular/http';

import { PathService } from "./path.service";
import { Collection } from "../models/collection";
import { RenditionConfiguration } from "../models/rendition-configuration";
import { Photo } from "../models/photo";

@Injectable()
export class PhotoService {

  constructor(
    private pathService: PathService,
    private http: Http
  ) { }

  recentPhotos(collection: Collection, renditionConfigurations: RenditionConfiguration[]): Promise<Photo[]> {
    let queryString = "";
    if (renditionConfigurations.length > 0) {
      queryString = `?rendition-configuration-ids=${renditionConfigurations.map((c => c.id)).join(",")}`;
    }
    let url = `${this.pathService.recentPhotos(collection)}${queryString}`
    return this.http
      .get(url)
      .toPromise()
      .then((response) => {
        let r = response.json() as PaginatedPhotos;

        r.data.map((photo) => {
          photo.collection = collection;
          return photo
        })
        return r.data;
      })
      .catch((e) => {
        return Promise.reject(e)
      });
  }
}

interface PaginatedPhotos {
  data: Photo[];
}
