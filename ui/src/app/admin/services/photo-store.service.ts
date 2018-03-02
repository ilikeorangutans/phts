import { RenditionConfiguration } from './../models/rendition-configuration';
import { Photo } from './../models/photo';
import { Observable } from 'rxjs/Observable';
import { Collection } from './../models/collection';
import { PhotoService } from './photo.service';
import { Injectable } from '@angular/core';

@Injectable()
export class PhotoStoreService {

  constructor(
    private photoService: PhotoService
  ) { }

  recentPhotos(collection: Observable<Collection>, renditionConfigurations: Array<RenditionConfiguration>): Observable<Array<Photo>> {
    console.log('PhotoStoreService.recentPhotos()');
    return collection.switchMap(c => {
      console.log('PhotoStoreService.recentPhotos() fetching recent for ', c.id);
      return this.photoService.recentPhotos(c, c.renditionConfigurations);
    });
  }



}
