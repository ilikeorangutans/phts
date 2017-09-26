import { Injectable, Inject } from '@angular/core';
import { DOCUMENT } from "@angular/common";

import { Collection } from "../models/collection";
import { Rendition } from "../models/rendition";

@Injectable()
export class PathService {

  constructor(
    @Inject(DOCUMENT) private document: any,
  ) { }


  apiBase(): string {
    this.document.location.href;

    return new URL("/admin/api/", "http://localhost:8080").toString();
  }

  collections(): string {
    return new URL("collections", this.apiBase()).toString()
  }

  showCollection(collection: Collection): string {
    return new URL(collection.slug, `${this.collections()}/`).toString()
  }

  recentPhotos(collection: Collection): string {
    let p = this.showCollection(collection);
    return new URL("photos/recent", `${p}/`).toString();
  }

  rendition(collection: Collection, rendition: Rendition): string {
    return new URL(`photos/renditions/${rendition.id}`, `${this.showCollection(collection)}/`).toString();
  }

  renditionConfigurations(collection: Collection): string {
    return new URL("rendition_configurations", `${this.showCollection(collection)}/`).toString();
  }

}
