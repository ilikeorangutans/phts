import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from "@angular/router";

import 'rxjs/add/operator/switchMap';

import { CollectionService } from "../../services/collection.service";
import { CurrentCollectionService } from "../current-collection.service";
import { RenditionConfigurationService } from "../../services/rendition-configuration.service";

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
    console.log("CollectionComponent::ngOnInit()");
    this.activatedRoute
      .paramMap
      .switchMap((params: ParamMap) => this.collectionService.bySlug(params.get('slug')))
      .subscribe(collection => {
        this.currentCollectionService.setCurrent(collection);

        this.renditionConfigService
          .forCollection(collection)
          .then(configs => this.currentCollectionService.setConfigs(configs));
      });
  }

}
