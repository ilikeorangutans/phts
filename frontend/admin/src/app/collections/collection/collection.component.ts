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

  constructor(
    private collectionService: CollectionService,
    private currentCollectionService: CurrentCollectionService,
    private activatedRoute: ActivatedRoute,
    private renditionConfigService: RenditionConfigurationService
  ) { }

  ngOnInit() {
    console.log('CollectionComponent::ngOnInit()');
    this.activatedRoute
      .paramMap
      .switchMap((params: ParamMap) => {
        console.log('CollectionComponent::ngOnInit() switchMap callback');
        return this.collectionService.bySlug(params.get('slug'));
      })
      .subscribe(
        (collection) => {
          this.renditionConfigService
            .forCollection(collection)
            .then(configs => this.registerCurrentCollection(collection, configs));
        }
      );
  }

  registerCurrentCollection(collection: Collection, configs: Array<RenditionConfiguration>) {
    collection.renditionConfigurations = configs;
    this.currentCollectionService.setCurrent(collection);
  }
}
