import { Photo } from './../../models/photo';
import { Collection } from './../../models/collection';
import { Paginator } from './../../models/paginator';
import { PathService } from './../../services/path.service';
import { PhotoService } from './../../services/photo.service';
import { CollectionService } from './../../services/collection.service';
import { Component, OnInit } from '@angular/core';
import { RenditionConfigurationService } from '../../services/rendition-configuration.service';
import { Rendition } from '../../models/rendition';

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

  constructor(
    private collectionService: CollectionService,
    private renditionConfigurationService: RenditionConfigurationService,
    private photoService: PhotoService,
    private pathService: PathService
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
          this.adminPreviewConfigID = configs.find(rc => rc.name === 'admin thumbnails').id;
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

  adminPreview(photo: Photo): Rendition {
    return photo.renditions.find(r => r.renditionConfigurationID === this.adminPreviewConfigID);
  }

  renditionURL(rendition): String {
    return this.pathService.rendition(this.collection, rendition);
  }

  loadMore(lastID: number, lastUpdatedAt: Date) {
    this.paginator = Paginator.newTimestampPaginator('updated_at', lastUpdatedAt);
    this.loadPhotos();
  }
}