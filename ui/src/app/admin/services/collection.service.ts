import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

import { PathService } from './path.service';
import { Collection } from '../models/collection';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class CollectionService {

  private currentCollection: BehaviorSubject<Collection> = new BehaviorSubject<Collection>(null);

  current: Observable<Collection> = this.currentCollection.asObservable();

  constructor(
    private http: HttpClient,
    private pathService: PathService
  ) { }

  setCurrent(collection: Collection) {
    if (this.currentCollection.getValue() === collection) {
      console.log('not setting current collection as it is the same');
      return;
    }
    this.currentCollection.next(collection);
  }

  bySlug(slug: string): Promise<Collection> {
    const url = this.pathService.collectionBase(slug);
    return this.http
      .get<Collection>(url)
      .toPromise()
      .then((c) => {
        c.createdAt = new Date(c.createdAt);
        c.updatedAt = new Date(c.updatedAt);

        return c;
      })
      .catch((e) => {
        console.log(e);
        return Promise.reject(e);
      });
  }

  recent(): Promise<Array<Collection>> {
    const url = this.pathService.collections();

    return this.http
      .get<Array<Collection>>(url)
      .toPromise()
      .then((collections) => {
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

  save(collection: Collection): Promise<Collection> {
    const url = this.pathService.collections();

    return this.http.post<Collection>(url, collection)
      .toPromise();
  }
}
