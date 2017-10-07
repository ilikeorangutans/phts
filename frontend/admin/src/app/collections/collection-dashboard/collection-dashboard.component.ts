import { Component, OnInit } from '@angular/core';

import { CurrentCollectionService } from "../current-collection.service";
import { PhotoService } from "../../services/photo.service";
import { PathService } from "../../services/path.service";
import { Collection } from "../../models/collection";
import { RenditionConfiguration } from "../../models/rendition-configuration";
import { Rendition } from "../../models/rendition";
import { Photo } from "../../models/photo";

@Component({
  selector: 'app-collection-dashboard',
  templateUrl: './collection-dashboard.component.html',
  styleUrls: ['./collection-dashboard.component.css']
})
export class CollectionDashboardComponent implements OnInit {

  collection: Collection = null;
  photos: Array<Photo> = new Array<Photo>();

  constructor(
    private currentCollectionService: CurrentCollectionService,
    private photoService: PhotoService,
    private pathService: PathService
  ) {
    currentCollectionService.current$.subscribe(collection => this.loadCollection(collection));

    currentCollectionService.renditionConfigs$.subscribe(configs => this.loadRecentPhotos(configs))
  }

  ngOnInit() {
  }

  loadCollection(collection: Collection) {
    this.collection = collection;
  }

  loadRecentPhotos(configs: Array<RenditionConfiguration>) {
    console.log("CollectionDashboardComponent::loadRecentPhotos()", configs)

    this.photoService
      .recentPhotos(this.collection, configs.filter(c => c.name == "admin thumbnails"))
      .then(photos => this.photos = photos);
  }

  renditionURI(rendition: Rendition): String {
    return this.pathService.rendition(this.collection, rendition);
  }

}