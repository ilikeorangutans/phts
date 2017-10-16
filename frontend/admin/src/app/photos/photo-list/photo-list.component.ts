import { PathService } from './../../services/path.service';
import { Rendition } from './../../models/rendition';
import { Paginator, TimestampPaginatorType } from './../../models/paginator';
import { Photo } from './../../models/photo';
import { Collection } from './../../models/collection';
import { PhotoService } from './../../services/photo.service';
import { CurrentCollectionService } from './../../collections/current-collection.service';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-photo-list',
  templateUrl: './photo-list.component.html',
  styleUrls: ['./photo-list.component.css']
})
export class PhotoListComponent implements OnInit {

  collection: Collection;

  photos: Array<Photo>;

  paginator: Paginator;

  adminPreviewConfigID: number;

  constructor(
    private currentCollectionService: CurrentCollectionService,
    private photoService: PhotoService,
    private pathService: PathService
  ) {
    this.paginator = Paginator.newTimestampPaginator('updated_at');
    this.currentCollectionService.current$.subscribe((c) => {
      if (c === null) {
        return;
      }

      this.collection = c;
      this.adminPreviewConfigID = this.collection.renditionConfigurations.find(rc => rc.name === 'admin thumbnails').id;

      this.loadPhotos();
    });
  }

  loadPhotos() {
    this.photoService.list(this.collection, this.paginator)
      .then(photos => {
        this.photos = photos;
      });
  }

  ngOnInit() {
  }

  adminPreview(photo: Photo): Rendition {
    return photo.renditions.find(r => r.renditionConfigurationID === this.adminPreviewConfigID);
  }

  renditionURL(rendition): String {
    return this.pathService.rendition(this.collection, rendition);
  }

}
