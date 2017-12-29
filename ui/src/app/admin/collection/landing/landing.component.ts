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
export class LandingComponent implements OnInit, OnDestroy {

  photos: Array<Photo> = new Array<Photo>();
  collection: Collection = null;

  private sub: Subscription;

  constructor(
    private photoService: PhotoService,
    private pathService: PathService,
    private collectionService: CollectionService,
    private renditionConfigurationService: RenditionConfigurationService
  ) { }

  ngOnInit() {
    this.sub = this.collectionService.current.subscribe(collection => {
      this.collection = collection;
      if (collection !== null) {
        this.loadRecentPhotos();
      }
    });
  }

  loadRecentPhotos() {
    this.renditionConfigurationService
      .forCollection(this.collection)
      .then(configs => {
        this.photoService
          .recentPhotos(this.collection, configs.filter(c => c.name === 'admin thumbnails'))
            .then(photos => this.photos = photos);
      });
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  renditionURI(rendition: Rendition): String {
    return this.pathService.rendition(this.collection, rendition);
  }

}
