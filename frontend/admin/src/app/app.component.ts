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

  ngOnInit(): void {
  }
}
