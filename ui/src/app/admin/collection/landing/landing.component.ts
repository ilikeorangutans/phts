import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/switchMap';

import { RenditionConfiguration } from './../../models/rendition-configuration';
import { CollectionService } from './../../services/collection.service';
import { PhotoStore } from './../../stores/photo.store';
import { Photo } from './../../models/photo';
import { PathService } from './../../services/path.service';
import { PhotoService } from './../../services/photo.service';
import { CollectionStore } from './../../stores/collection.store';
import { Collection } from '../../models/collection';
import { Rendition } from '../../models/rendition';

@Component({
  selector: 'app-landing',
  templateUrl: './landing.component.html',
  styleUrls: ['./landing.component.css'],
  providers: [PhotoStore]
})
export class LandingComponent implements OnInit {

  photos: Observable<Array<Photo>>;
  collection: Collection = null;
  previewRendition: RenditionConfiguration;

  constructor(
    private collectionService: CollectionService,
    private collectionStore: CollectionStore,
    private photoService: PhotoService,
    private pathService: PathService,
    private photoStore: PhotoStore,
    private router: Router
  ) { }

  ngOnInit(): void {
    this.collectionStore.current.subscribe(c => {
      this.collection = c;
      this.previewRendition = c.renditionConfigurations.find(r => r.name === 'admin thumbnails');
    });
    this.photos = this.photoStore.recent;
    this.refreshRecentPhotos();
  }

  refreshRecentPhotos(): void {
    this.photoStore.refreshRecent();
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
