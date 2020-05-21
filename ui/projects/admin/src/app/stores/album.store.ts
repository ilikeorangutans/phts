import { Injectable } from '@angular/core';

import { Collection } from './../models/collection';
import { PhotoService } from './../services/photo.service';
import { Paginator } from './../models/paginator';
import { Photo } from './../models/photo';
import { Album } from './../models/album';
import { BehaviorSubject, Observable } from 'rxjs';
import { first } from 'rxjs/operators';

@Injectable()
export class AlbumStore {
  private readonly _list: BehaviorSubject<Array<Photo>> = new BehaviorSubject<
    Array<Photo>
  >([]);

  readonly list: Observable<Array<Photo>> = this._list.asObservable();

  constructor(
    private readonly collection: Collection,
    private readonly album: Album,
    private readonly photoService: PhotoService
  ) {}

  loadPhotos(paginator: Paginator): void {
    this.photoService
      .listAlbum(this.collection, this.album, paginator)
      .pipe(first())
      .subscribe((photos) => this._list.next(photos));
  }
}
