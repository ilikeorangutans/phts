import { Component, OnInit } from '@angular/core';

import { Collection } from '../../models/collection';
import { Observable } from 'rxjs';

import { RenditionConfiguration } from './../../models/rendition-configuration';
import { PhotoStore } from './../../stores/photo.store';
import { Photo } from './../../models/photo';
import { CollectionStore } from './../../stores/collection.store';

@Component({
  selector: 'app-landing',
  templateUrl: './landing.component.html',
  styleUrls: ['./landing.component.css'],
})
export class LandingComponent implements OnInit {
  photos: Observable<Array<Photo>>;
  collection: Collection | null = null;
  previewRendition: RenditionConfiguration;

  constructor(
    private collectionStore: CollectionStore,
    private photoStore: PhotoStore
  ) {}

  ngOnInit(): void {
    this.collectionStore.current.subscribe((c) => {
      if (c === null) {
        throw 'collection is null';
      }
      this.collection = c;

      const previewRendition = c.renditionConfigurations.find(
        (r) => r.name === 'admin thumbnails'
      );
      if (previewRendition === undefined) {
        throw 'previewRendition is undefined';
      }

      this.previewRendition = previewRendition;
    });
    this.photos = this.photoStore.recent;
    this.refreshRecentPhotos();
  }

  refreshRecentPhotos(): void {
    this.photoStore.refreshRecent();
  }

  onPhotoClicked(photo: Photo): void {
    alert(`show preview for photo ${photo.id}`);
  }
}
