import { Injectable, Inject } from '@angular/core';
import { DOCUMENT } from "@angular/common";

import { Collection } from "../models/collection";

const apiBase: string = "/admin/api/";

@Injectable()
export class PathService {

  constructor(
    @Inject(DOCUMENT) private document: any,
  ) { }


  apiBase(): URL {
    this.document.location.href;

    return new URL(apiBase, "http://localhost:8080");
  }

  collections(): URL {
    return new URL("collections", this.apiBase().toString())
  }

  showCollection(collection: Collection): URL {
    return new URL(collection.id.toString(), this.collections.toString())
  }

}
