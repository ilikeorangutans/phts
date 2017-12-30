import { RenditionConfiguration } from './../../models/rendition-configuration';
import { Photo } from './../../models/photo';
import { Collection } from './../../models/collection';
import { ActivatedRoute } from '@angular/router';
import { PathService } from './../../services/path.service';
import { PhotoService } from './../../services/photo.service';
import { Component, OnInit } from '@angular/core';
import { CollectionService } from '../../services/collection.service';
import { Subscription } from 'rxjs/Subscription';
import { OnDestroy } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import 'rxjs/add/operator/filter';
import { Rendition } from '../../models/rendition';
import { RenditionConfigurationService } from '../../services/rendition-configuration.service';

@Component({
  selector: 'app-photo-details',
  templateUrl: './photo-details.component.html',
  styleUrls: ['./photo-details.component.css']
})
export class PhotoDetailsComponent implements OnInit, OnDestroy {
  private sub: Subscription;

  photo: Photo;
  collection: Collection;
  configs: Array<RenditionConfiguration>;
  previewID = 0;
  adminPreviewConfigID = 0;

  constructor(
    private collectionService: CollectionService,
    private photoService: PhotoService,
    private activatedRoute: ActivatedRoute,
    private pathService: PathService,
    private renditionConfigurationService: RenditionConfigurationService
  ) { }

  ngOnInit() {

    this.sub = this.collectionService.current
    .filter(collection => collection !== null)
    .switchMap(collection => {
      this.collection = collection;
      return this.renditionConfigurationService.forCollection(collection);
    }).switchMap(configs => {
      this.configs = configs;
      this.adminPreviewConfigID = configs.find(rc => rc.name === 'admin preview').id;
      return this.activatedRoute.params.map(params => +params['photo_id']);
    }).switchMap(photoID => {
      return this.photoService.byID(this.collection, photoID, []);
    }).subscribe(photo => {
      this.previewID = photo.renditions.find(r => r.renditionConfigurationID === this.adminPreviewConfigID).id;
      this.photo = photo;
    });
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  configByRendition(rendition: Rendition): RenditionConfiguration {
    return this.configs.find((c) => c.id === rendition.renditionConfigurationID);
  }

  preview(): Rendition {
    return this.photo.renditions.find(r => r.id === this.previewID);
  }

  renditionURL(rendition): String {
    return this.pathService.rendition(this.collection, rendition);
  }

  selectPreview(rendition: Rendition) {
    this.previewID = rendition.id;
  }
}
