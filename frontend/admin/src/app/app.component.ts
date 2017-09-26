import { Component, Inject, OnInit } from '@angular/core';
import { CollectionService } from "./services/collection.service";
import { PhotoService } from "./services/photo.service";
import { PathService } from "./services/path.service";
import { RenditionConfigurationService } from "./services/rendition-configuration.service";
import { Collection } from "./models/collection";
import { Photo } from "./models/photo";
import { RenditionConfiguration } from "./models/rendition-configuration";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  constructor(
    private pathService: PathService,
    private collectionService: CollectionService,
    private photoService: PhotoService,
    private renditionConfigurationService: RenditionConfigurationService
  ) {}

  title = 'app';

  collections: Collection[];
  photos: Photo[];
  configurations: RenditionConfiguration[];

  ngOnInit(): void {
    console.log("AppComponent::ngOnInit()");
    this.collectionService.recent().then((c) => {
      this.collections = c;

      console.log("Got collections");
      console.log(this.collections);
      let collection = this.collections[0];

      this.renditionConfigurationService.forCollection(collection).then((c) => {
        console.log(c);
        this.configurations = c;

        let adminConfigs = c.filter((c) => c.name.startsWith("admin"))

        this.photoService.recentPhotos(collection, adminConfigs).then((p) => {
          this.photos = p;
        });
      });
    });


  }

}
