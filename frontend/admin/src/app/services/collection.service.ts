import { Injectable } from '@angular/core';
import { Http } from '@angular/http';

import 'rxjs/add/operator/toPromise';

import { Collection } from '../models/collection';
import { PathService } from './path.service';

@Injectable()
export class CollectionService {

  constructor(
    private http: Http,
    private pathService: PathService
  ) { }

  bySlug(slug: string): Promise<Collection> {
    const url = this.pathService.collectionBase(slug);
    console.log('CollectionService::bySlug()', url);
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

  recent(): Promise<Array<Collection>> {
    const url = this.pathService.collections();
    console.log('CollectionServivce::recent() from', url);

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

  save(collection: Collection): Promise<Collection> {
    const url = this.pathService.collections();
    console.log('CollectionServivce::save() to', url)

    return this.http.post(url, collection)
      .toPromise()
      .then((response) => {
        console.log(response);
        return response.json() as Collection;
      })
      .catch((e) => {
        console.log(e);
        return null;
      });

  }
}
