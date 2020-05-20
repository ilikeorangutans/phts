import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { PathService } from './path.service';
import { Collection } from '../models/collection';

import 'rxjs/add/operator/distinctUntilChanged';
import 'rxjs/add/operator/filter';
import { BehaviorSubject, Observable } from 'rxjs';
import { filter, distinctUntilChanged, first, map } from 'rxjs/operators';

type CollectionOrNull = Collection | null;

@Injectable()
export class CollectionService {
  private currentCollection: BehaviorSubject<
    CollectionOrNull
  > = new BehaviorSubject<CollectionOrNull>(null);

  current: Observable<CollectionOrNull> = this.currentCollection.pipe(
    filter((c) => c !== null),
    distinctUntilChanged()
  );

  constructor(private http: HttpClient, private pathService: PathService) {}

  setCurrent(collection: Collection) {
    if (this.currentCollection.getValue() === collection) {
      return;
    }
    this.currentCollection.next(collection);
  }

  bySlug(slug: string): Observable<Collection> {
    const url = this.pathService.collection(slug);
    return this.http.get<Collection>(url).pipe(
      map((c) => {
        c.createdAt = new Date(c.createdAt);
        c.updatedAt = new Date(c.updatedAt);

        return c;
      }),
      first()
    );
  }

  recent(): Observable<Array<Collection>> {
    const url = this.pathService.collections();

    return this.http.get<Array<Collection>>(url).pipe(
      map((collections) => {
        collections = collections.map((c) => {
          c.createdAt = new Date(c.createdAt);
          c.updatedAt = new Date(c.updatedAt);
          return c;
        });

        return collections;
      })
    );
  }

  save(collection: Collection): Observable<Collection> {
    const url = this.pathService.collections();

    return this.http.post<Collection>(url, collection);
  }

  delete(collection: Collection) {
    const url = this.pathService.collection(collection.slug);
    console.log(url);
    this.http.delete(url).subscribe((response) => {
      console.log('response is ', response);
    });
  }
}
