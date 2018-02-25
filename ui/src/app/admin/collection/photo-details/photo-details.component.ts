import { Component, OnInit } from '@angular/core';
import { OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { Subscription } from 'rxjs/Subscription';

import { Collection } from './../../models/collection';
import { Photo } from './../../models/photo';
import { Rendition } from '../../models/rendition';
import { RenditionConfiguration } from './../../models/rendition-configuration';
import { PhotoService } from './../../services/photo.service';

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

  previewRendition: RenditionConfiguration;

  constructor(
    private photoService: PhotoService,
    private activatedRoute: ActivatedRoute,
    private router: Router
  ) { }

  setCollection(collection: Collection) {
    this.collection = collection;
    this.configs = collection.renditionConfigurations;
    this.previewRendition = this.configs.find(rc => rc.name === 'admin preview');

    this.sub = this.activatedRoute.params
      .map(params => +params['photo_id'])
      .switchMap(photoID => this.photoService.byID(this.collection, photoID, []))
      .subscribe(photo => {
        this.photo = photo;
      });
  }

  ngOnInit() {
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  configByRendition(rendition: Rendition): RenditionConfiguration {
    return this.configs.find((c) => c.id === rendition.renditionConfigurationID);
  }

  selectPreview(configID: number) {
    this.previewRendition = this.configs.find(r => r.id === configID);
  }

  delete(): void {
    this.photoService.delete(this.collection, this.photo);
    this.router.navigate(['admin', 'collection', this.collection.slug]);
  }
}
