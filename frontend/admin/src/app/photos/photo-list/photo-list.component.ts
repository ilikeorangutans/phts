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

  constructor(
    private currentCollectionService: CurrentCollectionService,
    private photoService: PhotoService
  ) {
    this.currentCollectionService.current$.subscribe((c) => {
      if (c === null) {
        return;
      }

      this.collection = c;
      this.loadPhotos();
    });
  }

  loadPhotos() {
    this.photoService.list(this.collection)
      .then(photos => {
        this.photos = photos;
      });
  }

  ngOnInit() {
  }

}
