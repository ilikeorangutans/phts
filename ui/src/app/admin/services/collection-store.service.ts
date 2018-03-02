import { Photo } from './../../public/models/photo';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Observable';
import { Collection } from './../models/collection';
import { CollectionService } from './collection.service';
import { Injectable } from '@angular/core';
import 'rxjs/add/operator/first';

@Injectable()
export class CollectionStoreService {
  private readonly _current: BehaviorSubject<Collection> = new BehaviorSubject<Collection>(null);

  private readonly _recent: Subject<Array<Collection>> = new Subject();

  readonly current: Observable<Collection> = this._current.filter(
    c => c !== null
  );

  constructor(private collectionService: CollectionService) {}

  recent(): Observable<Array<Collection>> {
    return this._recent.asObservable();
  }

  setCurrent(collection: Collection) {
    if (this._current.getValue() === collection) {
      return;
    }
    this._current.next(collection);
    this._recent.next([]);
  }

  save(collection: Collection) {
    this.collectionService
      .save(collection)
      .first()
      .subscribe(c => {
        this.refreshRecent();
      });
  }

  refreshRecent() {
    console.log('CollectionStoreService.refreshRecent()');
    this.collectionService
      .recent()
      .first()
      .subscribe(collections => this._recent.next(collections));
  }
}
