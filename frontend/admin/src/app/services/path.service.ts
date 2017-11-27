import { Injectable, Inject } from '@angular/core';
import { DOCUMENT } from '@angular/common';

import { Collection } from '../models/collection';
import { Rendition } from '../models/rendition';

@Injectable()
export class PathService {

  constructor(
    @Inject(DOCUMENT) private document: any,
  ) { }


  apiBase(): string {
    return new URL('/admin/api/', 'http://localhost:8080').toString();
  }

  collections(): string {
    return new URL('collections', this.apiBase()).toString();
  }

  collectionBase(slug: string): string {
    return new URL(slug, `${this.collections()}/`).toString();
  }

  uploadPhoto(collection: Collection): string {
    const p = this.collectionBase(collection.slug);
    return new URL('photos', `${p}/`).toString();
  }

  recentPhotos(collection: Collection): string {
    const p = this.collectionBase(collection.slug);
    return new URL('photos/recent', `${p}/`).toString();
  }

  rendition(collection: Collection, rendition: Rendition): string {
    return new URL(`photos/renditions/${rendition.id}`, `${this.collectionBase(collection.slug)}/`).toString();
  }

  renditionConfigurations(collection: Collection): string {
    return new URL('rendition_configurations', `${this.collectionBase(collection.slug)}/`).toString();
  }

  showPhoto(collection: Collection, photoID: number): string {
    return new URL(`photos/${photoID}`, `${this.collectionBase(collection.slug)}/`).toString();
  }

  listPhotos(collection: Collection): string {
    return new URL(`photos`, `${this.collectionBase(collection.slug)}/`).toString();
  }

  shareSites(): string {
    return new URL('share-sites', this.apiBase()).toString();
  }

}
