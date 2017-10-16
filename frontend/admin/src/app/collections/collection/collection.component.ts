import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';

import 'rxjs/add/operator/switchMap';

import { Collection } from '../../models/collection';
import { RenditionConfiguration } from '../../models/rendition-configuration';
import { CollectionService } from '../../services/collection.service';
import { CurrentCollectionService } from '../current-collection.service';
import { RenditionConfigurationService } from '../../services/rendition-configuration.service';

@Component({
  selector: 'app-collection',
  templateUrl: './collection.component.html',
  styleUrls: ['./collection.component.css'],
  providers: [CurrentCollectionService]
})
export class CollectionComponent implements OnInit {

  collection: Collection;

  constructor(
    private collectionService: CollectionService,
    private currentCollectionService: CurrentCollectionService,
    private activatedRoute: ActivatedRoute,
    private renditionConfigService: RenditionConfigurationService
  ) { }

  ngOnInit() {
    this.activatedRoute
      .paramMap
      .switchMap((params: ParamMap) => {
        return this.collectionService.bySlug(params.get('slug'));
      })
      .subscribe(
        (collection) => {
          this.collection = collection;
          this.renditionConfigService
            .forCollection(collection)
            .then((configs) => {
              this.registerCurrentCollection(this.collection, configs);
            });
        }
      );
  }

  registerCurrentCollection(collection: Collection, configs: Array<RenditionConfiguration>) {
    collection.renditionConfigurations = configs;
    this.currentCollectionService.setCurrent(collection);
  }
}
