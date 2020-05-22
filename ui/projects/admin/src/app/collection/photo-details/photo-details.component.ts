import { Component, OnInit } from '@angular/core';
import { OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { Observable } from 'rxjs';
import { map, switchMap, first } from 'rxjs/operators';

import { PhotoShares, ShareService } from './../../services/share.service';
import { CollectionStore } from './../../stores/collection.store';
import { Collection } from './../../models/collection';
import { Photo } from './../../models/photo';
import { Rendition } from '../../models/rendition';
import { RenditionConfiguration } from './../../models/rendition-configuration';
import { PhotoService } from './../../services/photo.service';

@Component({
  selector: 'app-photo-details',
  templateUrl: './photo-details.component.html',
  styleUrls: ['./photo-details.component.css'],
})
export class PhotoDetailsComponent implements OnInit, OnDestroy {
  photo: Observable<Photo>;
  collection: Collection;
  configs: Array<RenditionConfiguration>;

  previewRendition: RenditionConfiguration;
  shares: PhotoShares;

  constructor(
    private collectionStore: CollectionStore,
    private photoService: PhotoService,
    private activatedRoute: ActivatedRoute,
    private router: Router,
    private shareService: ShareService
  ) {}

  setCollection(collection: Collection) {
    this.collection = collection;
    this.configs = collection.renditionConfigurations;
    const previewRendition = this.configs.find(
      (rc) => rc.name === 'admin preview'
    );
    if (previewRendition === undefined) {
      throw 'previewRendition  is undefined';
    }
    this.previewRendition = previewRendition;

    this.photo = this.activatedRoute.params.pipe(
      map((params) => +params['photo_id']),
      switchMap((photoID) =>
        this.photoService.byID(this.collection, photoID, [])
      )
    );

    this.photo
      .pipe(first())
      .subscribe(
        (photo) =>
          (this.shares = new PhotoShares(
            this.shareService,
            this.collection,
            photo
          ))
      );
  }

  ngOnInit() {
    this.collectionStore.current.pipe(first()).subscribe((collection) => {
      if (collection === null) {
        throw 'collection is null';
      }
      this.setCollection(collection);
    });
  }

  ngOnDestroy(): void {}

  configByRendition(rendition: Rendition): RenditionConfiguration {
    const config = this.configs.find(
      (c) => c.id === rendition.renditionConfigurationID
    );

    if (config === undefined) {
      throw 'config is undefined';
    }
    return config;
  }

  selectPreview(configID: number) {
    const previewRendition = this.configs.find((r) => r.id === configID);
    if (previewRendition === undefined) {
      throw 'previewRendition is undefined';
    }
    this.previewRendition = previewRendition;
  }

  delete(photo: Photo): void {
    if (!confirm('Delete photo?')) {
      return;
    }
    this.photoService.delete(this.collection, photo);
    this.router.navigate(['collection', this.collection.slug]);
  }
}
