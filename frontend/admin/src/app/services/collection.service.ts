import { Injectable } from '@angular/core';
import { Http } from "@angular/http";

import 'rxjs/add/operator/toPromise';

import { Collection } from "../models/collection";
import { PathService } from "./path.service";

@Injectable()
export class CollectionService {

  constructor(
    private http: Http,
    private pathService: PathService
  ) { }

  recent(): Promise<Array<Collection>> {
    console.log("CollectionService::recent()");

    let url = this.pathService.collections().toString();
    console.log("Fetching collections from ", url);

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
