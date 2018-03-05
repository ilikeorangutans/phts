import { Paginator } from './../../models/paginator';
import { Component, OnInit } from '@angular/core';

import { Observable } from 'rxjs/Observable';

import { Photo } from './../../models/photo';
import { Collection } from './../../models/collection';
import { PhotoStore } from './../../stores/photo.store';
import { CollectionStore } from './../../stores/collection.store';
import { RenditionConfiguration } from '../../models/rendition-configuration';

@Component({
  selector: 'app-collection-photo-list',
  templateUrl: './collection-photo-list.component.html',
  styleUrls: ['./collection-photo-list.component.css']
})
export class CollectionPhotoListComponent implements OnInit {

  collection: Collection;
  photos: Observable<Array<Photo>>;
  paginator: Paginator;
  thumbnail: RenditionConfiguration;

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
    this.photos = this.photoStore.list;
    this.photoStore.updateList(this.paginator);
  }

  onPhotoClicked(photo: Photo): void {
    alert(`Show preview of photo ${photo.id}`);
  }

}
