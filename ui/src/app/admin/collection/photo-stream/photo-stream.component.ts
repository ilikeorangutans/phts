import { Component, OnInit, Output, EventEmitter } from '@angular/core';

import { SelectedPhoto } from './../selectable-photo-container/selectable-photo-container.component';
import { RenditionConfiguration } from './../../models/rendition-configuration';
import { Photo } from './../../models/photo';
import { Collection } from './../../models/collection';
import { Paginator } from './../../models/paginator';
import { PathService } from './../../services/path.service';
import { PhotoService } from './../../services/photo.service';
import { CollectionService } from './../../services/collection.service';
import { RenditionConfigurationService } from '../../services/rendition-configuration.service';
import { Rendition } from '../../models/rendition';
import { Album } from '../../models/album';
import { AlbumService } from '../../services/album.service';

@Component({
  selector: 'app-photo-stream',
  templateUrl: './photo-stream.component.html',
  styleUrls: ['./photo-stream.component.css']
})
export class PhotoStreamComponent implements OnInit {
  photos: Array<Photo> = [];
  paginator: Paginator;
  collection: Collection;
  adminPreviewConfigID: number;
  previewRenditionConfig: RenditionConfiguration;

  private selectedPhotos: Array<Photo> = [];

  constructor(
    private collectionService: CollectionService,
    private renditionConfigurationService: RenditionConfigurationService,
    private photoService: PhotoService,
    private pathService: PathService,
    private albumService: AlbumService
  ) { }

  ngOnInit() {
    this.paginator = Paginator.newTimestampPaginator('updated_at');
    this.collectionService.current.subscribe((c) => {
      if (c === null) {
        return;
      }

      this.collection = c;

      this.renditionConfigurationService.forCollection(c)
        .then(configs => {
          this.previewRenditionConfig = configs.find(rc => rc.name === 'admin thumbnails');
          this.adminPreviewConfigID = this.previewRenditionConfig.id;
          this.loadPhotos();
        });
    });
  }

  loadPhotos() {
    this.photoService.list(this.collection, this.paginator)
      .then(photos => {
        this.photos = this.photos.concat(photos);
      });
  }

  loadMore(lastID: number, lastUpdatedAt: Date) {
    this.paginator = Paginator.newTimestampPaginator('updated_at', lastUpdatedAt);
    this.loadPhotos();
  }

  shareSelectionToAlbum(album: Album) {
    this.albumService.addPhotos(this.collection, album, this.selectedPhotos);
  }

  setSelectedPhotos(photos: Array<Photo>) {
    this.selectedPhotos = photos;
  }
}
