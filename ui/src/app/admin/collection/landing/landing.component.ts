import { RenditionConfigurationService } from './../../services/rendition-configuration.service';
import { Photo } from './../../models/photo';
import { PathService } from './../../services/path.service';
import { PhotoService } from './../../services/photo.service';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import 'rxjs/add/operator/switchMap';

import { CollectionService } from './../../services/collection.service';
import { Collection } from '../../models/collection';
import { OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Rendition } from '../../models/rendition';

@Component({
  selector: 'app-landing',
  templateUrl: './landing.component.html',
  styleUrls: ['./landing.component.css']
})
export class LandingComponent {

  photos: Array<Photo> = new Array<Photo>();
  collection: Collection = null;

  constructor(
    private photoService: PhotoService,
    private pathService: PathService,
    private collectionService: CollectionService,
    private renditionConfigurationService: RenditionConfigurationService
  ) { }

  setCollection(collection: Collection) {
    this.collection = collection;
    this.loadRecentPhotos();
  }

  loadRecentPhotos() {
    this.photoService
      .recentPhotos(this.collection, this.collection.renditionConfigurations.filter(c => c.name === 'admin thumbnails'))
      .then(photos => this.photos = photos);
  }

  renditionURI(rendition: Rendition): String {
    return this.pathService.rendition(this.collection, rendition);
  }

}
