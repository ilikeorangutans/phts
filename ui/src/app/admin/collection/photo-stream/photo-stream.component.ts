import { Component, OnInit, Output, EventEmitter } from '@angular/core';

import { Observable } from 'rxjs/Observable';

import { CollectionStore } from './../../stores/collection.store';
import { PhotoStore } from '../../stores/photo.store';
import { RenditionConfiguration } from './../../models/rendition-configuration';
import { Photo } from './../../models/photo';
import { Collection } from './../../models/collection';
import { Paginator } from './../../models/paginator';
import { Album } from '../../models/album';
import { AlbumService } from '../../services/album.service';

@Component({
  selector: 'app-photo-stream',
  templateUrl: './photo-stream.component.html',
  styleUrls: ['./photo-stream.component.css']
})
export class PhotoStreamComponent implements OnInit {
  photos: Observable<Array<Photo>>;
  paginator: Paginator;
  collection: Collection;

  thumbnailRendition: RenditionConfiguration;

  private selectedPhotos: Array<Photo> = [];

  constructor(
    private collectionStore: CollectionStore,
    private photoStore: PhotoStore
  ) { }

  ngOnInit() {
    this.paginator = Paginator.newTimestampPaginator('updated_at');
    this.collectionStore.current
      .first()
      .subscribe(c => {
        this.collection = c;
        this.thumbnailRendition = c.renditionConfigurations.find(r => r.name === 'admin thumbnails');
      });

    this.photos = this.photoStore.list;
    this.refreshPhotos();
  }

  refreshPhotos(): void {
    this.photoStore.updateList(this.paginator);
  }

  loadMore(lastID: number, lastUpdatedAt: Date) {
    this.paginator = Paginator.newTimestampPaginator('updated_at', lastUpdatedAt);
    this.photoStore.updateList(this.paginator);
  }

  shareSelectionToAlbum(album: Album) {
    // this.albumService.addPhotos(this.collection, album, this.selectedPhotos);
  }

  setSelectedPhotos(photos: Array<Photo>) {
    this.selectedPhotos = photos;
  }
}
