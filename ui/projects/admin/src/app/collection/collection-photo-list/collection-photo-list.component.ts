import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';

import { OverlayComponent } from './../../shared/overlay/overlay.component';
import { Album } from './../../models/album';
import { Paginator } from './../../models/paginator';

import { Photo } from './../../models/photo';
import { Collection } from './../../models/collection';
import { PhotoStore } from './../../stores/photo.store';
import { CollectionStore } from './../../stores/collection.store';
import { RenditionConfiguration } from '../../models/rendition-configuration';
import { AlbumService } from '../../services/album.service';
import { Observable, Subscription, BehaviorSubject } from 'rxjs';
import { first, scan } from 'rxjs/operators';

@Component({
  selector: 'app-collection-photo-list',
  templateUrl: './collection-photo-list.component.html',
  styleUrls: ['./collection-photo-list.component.css'],
})
export class CollectionPhotoListComponent implements OnInit, OnDestroy {
  collection: Collection;
  photos: Observable<Array<Photo>>;
  paginator: Paginator;
  thumbnail: RenditionConfiguration;

  lastPhoto: Photo;

  showAlbumSelector = false;

  albums: Observable<Array<Album>>;
  private sub: Subscription;

  @ViewChild(OverlayComponent)
  overlay: OverlayComponent;

  photo: Photo;

  private readonly _organizePhotoMode: BehaviorSubject<
    boolean
  > = new BehaviorSubject(false);
  readonly organizePhotoMode: Observable<
    boolean
  > = this._organizePhotoMode.asObservable();

  constructor(
    private readonly collectionStore: CollectionStore,
    private readonly photoStore: PhotoStore,
    private readonly albumService: AlbumService
  ) {}

  ngOnInit() {
    this.paginator = Paginator.newTimestampPaginator('updated_at');
    this.collectionStore.current.pipe(first()).subscribe((c) => {
      if (c === null) {
        throw 'collection is null';
      }
      this.collection = c;
      const thumbnail = c.renditionConfigurations.find(
        (rc) => rc.name === 'admin thumbnails'
      );

      if (thumbnail === undefined) {
        throw 'thumbnail is undefined';
      }
      this.thumbnail = thumbnail;

      this.albums = this.albumService.list(c);
    });
    this.photos = this.photoStore.list.pipe(
      scan((acc, value) => {
        // Accumulate loaded photos in the acc array
        return acc.concat(value);
      }, new Array<Photo>())
    );
    this.photoStore.updateList(this.paginator);

    this.sub = this.photos.subscribe((photos) => {
      if (photos.length === 0) {
        return;
      }

      this.lastPhoto = photos[photos.length - 1];
    });
  }

  reloadPhotos(): void {
    this.photoStore.updateList(this.paginator);
  }

  onPhotoClicked(photo: Photo): void {
    this.photo = photo;
    this.overlay.show();
  }

  loadMore(): void {
    let ts = new Date();
    if (this.lastPhoto !== null) {
      ts = this.lastPhoto.updatedAt;
    }

    this.paginator = Paginator.newTimestampPaginator('updated_at', ts);

    this.reloadPhotos();
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  delete(photos: Array<Photo>): void {
    if (!confirm(`Delete ${photos.length} photos?`)) {
      return;
    }

    alert('TODO: implement delete');
  }

  showAlbumDialog(): void {
    this.showAlbumSelector = true;
  }

  shareSelectionToAlbum(album: Album, photos: Array<Photo>): void {
    this.albumService.addPhotos(this.collection, album, photos);
  }

  toggleOrganizePhotoMode(): void {
    this._organizePhotoMode.next(!this._organizePhotoMode.getValue());
  }
}
