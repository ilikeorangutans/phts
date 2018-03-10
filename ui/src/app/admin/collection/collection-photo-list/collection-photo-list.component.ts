import { Paginator } from './../../models/paginator';
import { Component, OnInit, OnDestroy } from '@angular/core';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/scan';

import { Photo } from './../../models/photo';
import { Collection } from './../../models/collection';
import { PhotoStore } from './../../stores/photo.store';
import { CollectionStore } from './../../stores/collection.store';
import { RenditionConfiguration } from '../../models/rendition-configuration';
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-collection-photo-list',
  templateUrl: './collection-photo-list.component.html',
  styleUrls: ['./collection-photo-list.component.css']
})
export class CollectionPhotoListComponent implements OnInit, OnDestroy {

  collection: Collection;
  photos: Observable<Array<Photo>>;
  paginator: Paginator;
  thumbnail: RenditionConfiguration;

  lastPhoto: Photo;

  private sub: Subscription;

  constructor(
    private collectionStore: CollectionStore,
    private photoStore: PhotoStore
  ) { }

  ngOnInit() {
    this.paginator = Paginator.newTimestampPaginator('updated_at');
    this.collectionStore.current.first()
      .subscribe(c => {
        this.collection = c;
        this.thumbnail = c.renditionConfigurations.find(rc => rc.name === 'admin thumbnails');
      });
    this.photos = this.photoStore.list.scan((acc, value) => {
      // Accumulate loaded photos in the acc array
      return acc.concat(value);
    }, []);
    this.photoStore.updateList(this.paginator);

    this.sub = this.photos.subscribe(photos => {
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
    alert(`Show preview of photo ${photo.id}`);
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

}
