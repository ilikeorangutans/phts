import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import { Subject } from "rxjs/Subject";

import 'rxjs/add/operator/toPromise';

import { Collection } from '../models/collection';
import { PathService } from './path.service';

@Injectable()
export class CollectionService {

  private currentCollectionSource = new Subject<Collection>();

  currentCollection$ = this.currentCollectionSource.asObservable();

  constructor(
    private http: Http,
    private pathService: PathService
  ) { }

  bySlug(slug: string): Promise<Collection> {
    const url = this.pathService.collectionBase(slug);
    console.log("CollectionService::bySlug()", url);
    return this.http
      .get(url)
      .toPromise()
      .then((response) => {
        const c = response.json() as Collection;

        c.createdAt = new Date(c.createdAt);
        c.updatedAt = new Date(c.updatedAt);

        return c;
      })
      .catch((e) => {
        console.log(e);
        return Promise.reject(e);
      });
  }

  setCurrent(collection: Collection) {
    console.log("Setting current ", collection)
    this.currentCollectionSource.next(collection);
  }

  recent(): Promise<Array<Collection>> {
    const url = this.pathService.collections();
    console.log("CollectionServivce::recent() from", url)

    return this.http
      .get(url)
      .toPromise()
      .then((response) => {
        let collections = response.json() as Collection[];

        collections = collections.map((c) => {
          c.createdAt = new Date(c.createdAt);
          c.updatedAt = new Date(c.updatedAt);
          return c;
        });

        return collections;
      })
      .catch((e) => {
        console.log(e);
        return Promise.reject(e);
      });
  }
}
