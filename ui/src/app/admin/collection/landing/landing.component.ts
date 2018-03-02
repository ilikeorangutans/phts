import { Observable } from 'rxjs/Observable';
import { PhotoStoreService } from './../../services/photo-store.service';
import { RenditionConfiguration } from './../../models/rendition-configuration';
import { RenditionConfigurationService } from './../../services/rendition-configuration.service';
import { Photo } from './../../models/photo';
import { PathService } from './../../services/path.service';
import { PhotoService } from './../../services/photo.service';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import 'rxjs/add/operator/switchMap';
import { from } from 'rxjs/observable/from';

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

  // photos: Array<Photo> = new Array<Photo>();
  photos: Observable<Array<Photo>>;
  collection: Collection = null;

  constructor(
    private photoStore: PhotoStoreService,
    private photoService: PhotoService,
    private pathService: PathService,
    private collectionService: CollectionService,
    private renditionConfigurationService: RenditionConfigurationService,
    private router: Router
  ) { }

  setCollection(collection: Collection) {
    this.collection = collection;

    this.photos = this.photoStore.recentPhotos(from([collection]), []);
  }

  previewRendition(): RenditionConfiguration {
    return this.collection.renditionConfigurations.find(c => c.name === 'admin thumbnails');
  }

  renditionURI(rendition: Rendition): String {
    return this.pathService.rendition(this.collection, rendition);
  }

  delete(): void {
    this.collectionService.delete(this.collection);

    alert('implement me: here we\'d delete this collection');

    this.router.navigate(['admin', 'collection']);
  }

}
