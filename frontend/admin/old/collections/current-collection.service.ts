import { Injectable } from '@angular/core';
import { Subject } from "rxjs/Subject";

import { Collection } from "../models/collection";

@Injectable()
export class CurrentCollectionService {

  private currentCollectionSource = new Subject<Collection>();

  current$ = this.currentCollectionSource.asObservable();

  constructor() { }

  setCurrent(collection: Collection) {
    console.log("CurrentCollectionService::setCurrent()", collection);
    this.currentCollectionSource.next(collection);
  }
}
