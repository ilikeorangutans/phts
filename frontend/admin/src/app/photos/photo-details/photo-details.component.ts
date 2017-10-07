import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from "@angular/router";

// import 'rxjs/add/operator/switchMap';

import { CurrentCollectionService } from "../../collections/current-collection.service";
import { PhotoService } from "../../services/photo.service";
import { PathService } from "../../services/path.service";
import { Collection } from "../../models/collection";
import { RenditionConfiguration } from "../../models/rendition-configuration";
import { Photo } from "../../models/photo";
import { Rendition } from "../../models/rendition";

@Component({
  selector: 'app-photo-details',
  templateUrl: './photo-details.component.html',
  styleUrls: ['./photo-details.component.css']
})
export class PhotoDetailsComponent implements OnInit {

  photo: Photo;
  collection: Collection;
  renditionConfigurations: Array<RenditionConfiguration>;

  photoID: number = 0;

  constructor(
    private currentCollectionService: CurrentCollectionService,
    private photoService: PhotoService,
    private activatedRoute: ActivatedRoute,
    private pathService: PathService
  ) {
    currentCollectionService.current$.subscribe(collection => {
      this.collection = collection;
      this.loadPhoto();
    });
    currentCollectionService.renditionConfigs$.subscribe(configs => {
      this.renditionConfigurations = configs;
      this.loadPhoto();
    });

    this.photoID = +activatedRoute.snapshot.params["photoID"];
  }

  loadPhoto() {
    console.log("PhotoDetailsComponent::loadPhoto()", this.collection)
    console.log("configs", this.renditionConfigurations);

    if (this.collection && this.renditionConfigurations && this.photoID) {
      this.photoService.byID(this.collection, this.photoID, this.renditionConfigurations)
        .then(photo => this.photo = photo);
    }
  }

  ngOnInit() {
  }

  configByRendition(rendition: Rendition): RenditionConfiguration {
    return this.renditionConfigurations.find((c) => c.id == rendition.renditionConfigurationID);
  }

  preview(): Rendition {
    let id = this.renditionConfigurations.find(rc => rc.name == "admin preview").id;
    return this.photo.renditions.find(r => r.renditionConfigurationID == id);
  }

  renditionURL(rendition): String {
    return this.pathService.rendition(this.collection, rendition);
  }
}
