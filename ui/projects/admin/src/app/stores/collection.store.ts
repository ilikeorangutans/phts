import { PhotoService } from './../services/photo.service';
import { Album } from './../models/album';
import { AlbumStore } from './album.store';
import { Collection } from './../models/collection';
import { CollectionService } from './../services/collection.service';
import { Injectable } from '@angular/core';
import { RenditionConfigurationService } from '../services/rendition-configuration.service';
import { BehaviorSubject, Observable } from 'rxjs';
import { filter, map, flatMap, first } from 'rxjs/operators';

@Injectable()
export class CollectionStore {
  private readonly _currentlyBusy: BehaviorSubject<
    boolean
  > = new BehaviorSubject<boolean>(false);

  private readonly _current: BehaviorSubject<Collection | null> = new BehaviorSubject<Collection | null>(
    null
  );

  private readonly _recent = new BehaviorSubject<Array<Collection>>([]);

  readonly current: Observable<Collection | null> = this._current.pipe(
    filter((c) => c !== null)
  );

  readonly recent: Observable<Array<Collection>> = this._recent.asObservable();

  readonly currentlyBusy: Observable<
    boolean
  > = this._currentlyBusy.asObservable();

  constructor(
    private readonly collectionService: CollectionService,
    private readonly renditionConfigurationService: RenditionConfigurationService,
    private readonly photoService: PhotoService
  ) {}

  private setBusy() {
    this._currentlyBusy.next(true);
  }

  private setIdle() {
    this._currentlyBusy.next(false);
  }

  setCurrentBySlug(slug: string): void {
    this.setBusy();
    this.collectionService
      .bySlug(slug)
      .pipe(
        flatMap((collection) => {
          return this.renditionConfigurationService
            .forCollection(collection)
            .pipe(
              map((configs) => {
                collection.renditionConfigurations = configs;
                return collection;
              })
            );
        })
      )
      .subscribe((collection) => {
        this._current.next(collection);
        this.setIdle();
      });
  }

  refreshRecent(): void {
    this.setBusy();
    this.collectionService
      .recent()
      .pipe(first())
      .subscribe((collections) => {
        this._recent.next(collections);
        this.setIdle();
      });
  }

  save(collection: Collection): void {
    this.setBusy();
    this.collectionService
      .save(collection)
      .pipe(first())
      .subscribe((_) => {
        this.refreshRecent();
      });
  }

  delete(collection: Collection): void {
    this.collectionService.delete(collection);
    this.refreshRecent();
  }

  albumStore(album: Album): AlbumStore {
    const collection = this._current.getValue();
    if (collection === null) {
      throw 'cant get album store if no collection selected';
    }
    return new AlbumStore(collection, album, this.photoService);
  }
}
