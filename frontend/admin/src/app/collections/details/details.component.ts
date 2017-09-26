import { Component, OnInit } from '@angular/core';
import 'rxjs/add/operator/switchMap';
import { ActivatedRoute, ParamMap, Params } from "@angular/router";

import { CollectionService } from "../../services/collection.service";
import { PhotoService } from "../../services/photo.service";
import { RenditionConfigurationService } from "../../services/rendition-configuration.service";
import { PathService } from "../../services/path.service";

import { Collection } from "../../models/collection";
import { Photo } from "../../models/photo";
import { RenditionConfiguration } from "../../models/rendition-configuration";

@Component({
  selector: 'app-details',
  templateUrl: './details.component.html',
  styleUrls: ['./details.component.css']
})
export class DetailsComponent implements OnInit {

  // TODO: might be better to check in the view if it's defined?
  collection: Collection = new Collection();

  photos: Photo[];

  configurations: RenditionConfiguration[];

  constructor(
    private activatedRoute: ActivatedRoute,
    private collectionService: CollectionService,
    private photoService: PhotoService,
    private renditionConfigService: RenditionConfigurationService,
    private pathService: PathService
  ) { }

  ngOnInit() {
    this.activatedRoute
      .paramMap
      .switchMap((params: ParamMap) => this.collectionService.bySlug(params.get("slug")))
      .subscribe(collection => this.loadCollection(collection));
  }

  loadCollection(collection: Collection) {
    console.log("DetailsComponent::loadCollection()");
    this.collection = collection;

    this.renditionConfigService
      .forCollection(collection)
      .then((configurations) => {
        this.configurations = configurations;

        let thumbnailConfig = this.configurations.filter(config => config.name == "admin thumbnails");

        this.photoService.recentPhotos(collection, thumbnailConfig)
          .then(photos => this.photos = photos);
      });

  }
}
