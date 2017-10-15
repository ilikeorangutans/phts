import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Collection } from '../models/collection';
import { RenditionConfiguration } from '../models/rendition-configuration';

@Injectable()
export class CurrentCollectionService {

  // TODO: this initial value being null messes things up.
  private currentCollectionSource: Subject<Collection> = new BehaviorSubject<Collection>(null);

  current$ = this.currentCollectionSource.asObservable();

  constructor() {}

  setCurrent(collection: Collection) {
    console.log('CurrentCollectionService::setCurrent()', collection);
    this.currentCollectionSource.next(collection);
  }
}
