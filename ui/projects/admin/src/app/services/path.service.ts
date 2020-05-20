import { Injectable } from '@angular/core';

import { BasePathService } from 'projects/shared/src/app/services/base-path.service';
import { Collection } from '../models/collection';
import { Rendition } from '../models/rendition';
import { Album } from '../models/album';

@Injectable({
  providedIn: 'root',
})
export class PathService {
  readonly api: string;

  constructor(private readonly basePath: BasePathService) {
    this.api = [this.basePath.apiHost, 'api', 'admin'].join('/');
  }

  joinToken(token: string): string {
    return [this.api, 'invite', token].join('/');
  }

  version(): string {
    return [this.api, 'version'].join('/');
  }

  authenticate(): string {
    return [this.api, 'authenticate'].join('/');
  }

  changePassword(): string {
    return [this.api, 'account', 'password'].join('/');
  }

  collections(): string {
    return [this.api, 'collections'].join('/');
  }

  collection(slug: string): string {
    return [this.collections(), slug].join('/');
  }

  uploadPhoto(collection: Collection): string {
    return [this.collection(collection.slug), 'photos'].join('/');
  }

  recentPhotos(collection: Collection): string {
    return [this.collection(collection.slug), 'photos/recent'].join('/');
  }

  rendition(collection: Collection, rendition: Rendition): string {
    return [
      this.collection(collection.slug),
      'photos/renditions',
      rendition.id,
    ].join('/');
  }

  renditionConfigurations(collection: Collection): string {
    return [this.collection(collection.slug), 'rendition_configurations'].join(
      '/'
    );
  }

  showPhoto(collection: Collection, photoID: number): string {
    return [this.collection(collection.slug), 'photos', photoID].join('/');
  }

  listPhotos(collection: Collection): string {
    return [this.collection(collection.slug), 'photos'].join('/');
  }

  shareSites(): string {
    return [this.api, 'share-sites'].join('/');
  }

  photoShares(collection: Collection, photoID: number): string {
    return [this.collection(collection.slug), 'photos', photoID, 'shares'].join(
      '/'
    );
  }

  albumBase(collection: Collection): string {
    return [this.collection(collection.slug), 'albums'].join('/');
  }

  albumDetails(collection: Collection, id: number): string {
    return [this.albumBase(collection), id].join('/');
  }

  albumPhotos(collection: Collection, album: Album): string {
    return [this.albumDetails(collection, album.id), 'photos'].join('/');
  }
}
