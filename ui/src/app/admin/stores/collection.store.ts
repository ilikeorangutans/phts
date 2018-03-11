import { PhotoService } from './../services/photo.service';
import { Album } from './../models/album';
import { AlbumStore } from './album.store';
import { Observable } from 'rxjs/Observable';
import { Collection } from './../models/collection';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { CollectionService } from './../services/collection.service';
import { Injectable } from '@angular/core';
import 'rxjs/add/operator/switchMap';
import 'rxjs/add/operator/mergeMap';
import 'rxjs/add/operator/first';
import { RenditionConfigurationService } from '../services/rendition-configuration.service';

@Injectable()
export class CollectionStore {

  private readonly _current: BehaviorSubject<Collection> = new BehaviorSubject<Collection>(null);

  private readonly _recent = new BehaviorSubject<Array<Collection>>([]);

  readonly current: Observable<Collection> = this._current.asObservable().filter(c => c !== null);

  readonly recent: Observable<Array<Collection>> = this._recent.asObservable();

  constructor(
    private readonly collectionService: CollectionService,
    private readonly renditionConfigurationService: RenditionConfigurationService,
    private readonly photoService: PhotoService
  ) { }

  setCurrentBySlug(slug: string): void {
    this.collectionService
      .bySlug(slug)
      .flatMap(collection => {
        return this.renditionConfigurationService
          .forCollection(collection)
          .map(configs => {
            collection.renditionConfigurations = configs;
            return collection;
          });
      })
      .subscribe(collection => {
        this._current.next(collection);
      });
  }

  refreshRecent(): void {
    this.collectionService.recent()
      .first()
      .subscribe(collections => {
        this._recent.next(collections);
      });
  }

  save(collection: Collection): void {
    this.collectionService
      .save(collection)
      .first()
      .subscribe(_ => {
        this.refreshRecent();
      });
  }

  delete(collection: Collection): void {
    this.collectionService
      .delete(collection);
    this.refreshRecent();
  }

  albumStore(album: Album): AlbumStore {
    return new AlbumStore(
      this._current.getValue(),
      album,
      this.photoService
    );
  }
}
