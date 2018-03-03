import { Photo } from './../models/photo';
import { PhotoService } from './../services/photo.service';
import { Observable } from 'rxjs/Observable';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Injectable } from '@angular/core';

import { CollectionStore } from './collection.store';

@Injectable()
export class PhotoStore {

  private readonly _recent = new BehaviorSubject<Array<Photo>>([]);

  readonly recent: Observable<Array<Photo>> = this._recent.asObservable();

  constructor(
    private collectionStore: CollectionStore,
    private photoService: PhotoService
  ) {}

  refreshRecent(): void {
    this.collectionStore.current
      .first()
      .switchMap(collection => {
        return this.photoService
          .recentPhotos(collection, []);
      })
      .subscribe(photos => this._recent.next(photos));
  }
}
