import { Album } from './../models/album';
import { Injectable, Inject } from '@angular/core';
import { DOCUMENT } from '@angular/common';

import { Collection } from '../models/collection';
import { Rendition } from '../models/rendition';
import { isDevMode } from '@angular/core';

@Injectable()
export class PathService {

  constructor(
    @Inject(DOCUMENT) private document: Document,
  ) { }

  apiHost(): string {
    if (isDevMode()) {
      return 'http://localhost:8080';
    } else {
      return this.document.baseURI;
    }
  }

  apiBase(): string {
    return new URL('/admin/api/', this.apiHost()).toString();
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

  photoShares(collection: Collection, photoID: number): string {
    return new URL(`photos/${photoID}/shares`, `${this.collectionBase(collection.slug)}/`).toString();
  }

  albumBase(collection: Collection): string {
    return [this.collectionBase(collection.slug), 'albums'].join('/');
  }

  albumDetails(collection: Collection, album: Album): string {
    return [this.albumBase(collection), album.id].join('/');
  }

  albumPhotos(collection: Collection, album: Album): string {
    return [this.albumDetails(collection, album), 'photos'].join('/');
  }

}
