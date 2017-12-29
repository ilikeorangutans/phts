import { RenditionConfiguration } from './../../models/rendition-configuration';
import { CollectionService } from './../../services/collection.service';
import { Component, OnInit } from '@angular/core';
import { RenditionConfigurationService } from '../../services/rendition-configuration.service';

import 'rxjs/add/operator/switchMap';

@Component({
  selector: 'app-rendition-configurations',
  templateUrl: './rendition-configurations.component.html',
  styleUrls: ['./rendition-configurations.component.css']
})
export class RenditionConfigurationsComponent implements OnInit {

  configurations: Array<RenditionConfiguration> = [];

  constructor(
    private collectionService: CollectionService,
    private renditionConfigurationService: RenditionConfigurationService
  ) { }

  ngOnInit() {
    this.collectionService.current.switchMap(collection => {
      if (collection === null) {
        return Promise.resolve([]);
      } else {
        return this.renditionConfigurationService.forCollection(collection);
      }
    })
    .subscribe(configs => this.configurations = configs);
  }

}
