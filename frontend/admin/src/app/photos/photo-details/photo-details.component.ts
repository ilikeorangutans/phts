import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';

// import 'rxjs/add/operator/switchMap';

import { CurrentCollectionService } from '../../collections/current-collection.service';
import { PhotoService } from '../../services/photo.service';
import { PathService } from '../../services/path.service';
import { Collection } from '../../models/collection';
import { RenditionConfiguration } from '../../models/rendition-configuration';
import { Photo } from '../../models/photo';
import { Rendition } from '../../models/rendition';

@Component({
  selector: 'app-photo-details',
  templateUrl: './photo-details.component.html',
  styleUrls: ['./photo-details.component.css']
})
export class PhotoDetailsComponent implements OnInit {

  photo: Photo;
  collection: Collection;

  photoID = 0;

  constructor(
    private currentCollectionService: CurrentCollectionService,
    private photoService: PhotoService,
    private activatedRoute: ActivatedRoute,
    private pathService: PathService
  ) {
    this.photoID = +activatedRoute.snapshot.params['photoID'];

    currentCollectionService.current$.subscribe(collection => {
      this.collection = collection;
      this.loadPhoto();
    });
  }

  loadPhoto() {
    console.log('PhotoDetailsComponent::loadPhoto()', this.collection);

    this.photoService.byID(this.collection, this.photoID, this.collection.renditionConfigurations)
      .then(photo => this.photo = photo);
  }

  ngOnInit() {
  }

  configByRendition(rendition: Rendition): RenditionConfiguration {
    return this.collection.renditionConfigurations.find((c) => c.id === rendition.renditionConfigurationID);
  }

  preview(): Rendition {
    const id = this.collection.renditionConfigurations.find(rc => rc.name === 'admin preview').id;
    return this.photo.renditions.find(r => r.renditionConfigurationID === id);
  }

  renditionURL(rendition): String {
    return this.pathService.rendition(this.collection, rendition);
  }
}
