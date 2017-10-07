import { Injectable } from '@angular/core';
import { Subject } from "rxjs/Subject";
import { Collection } from "../models/collection";
import { RenditionConfiguration } from "../models/rendition-configuration";

@Injectable()
export class CurrentCollectionService {

  private currentCollectionSource = new Subject<Collection>();
  private renditionConfigsSource = new Subject<RenditionConfiguration[]>();

  current$ = this.currentCollectionSource.asObservable();
  renditionConfigs$ = this.renditionConfigsSource.asObservable();

  constructor() {}

  setCurrent(collection: Collection) {
    this.currentCollectionSource.next(collection);
    this.renditionConfigsSource.next(new Array<RenditionConfiguration>());
  }

  setConfigs(configs: RenditionConfiguration[]) {
    this.renditionConfigsSource.next(configs);
  }
}
