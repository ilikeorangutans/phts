import { Observable } from 'rxjs/Observable';
import { RenditionConfigurationService } from './../../services/rendition-configuration.service';
import { ActivatedRoute } from '@angular/router';
import { CollectionService } from './../../services/collection.service';
import { Collection } from './../../models/collection';
import { Component, OnInit, EventEmitter, Output } from '@angular/core';
import 'rxjs/add/operator/merge';
import 'rxjs/add/operator/do';

@Component({
  selector: 'app-collection-loader',
  template: ''
})
export class CollectionLoaderComponent implements OnInit {

  @Output()
  collectionLoaded: EventEmitter<Collection> = new EventEmitter<Collection>();

  constructor(
    private collectionService: CollectionService,
    private renditionConfigurationService: RenditionConfigurationService,
    private route: ActivatedRoute
  ) { }

  ngOnInit() {
    this.collectionService
      .current
      .switchMap(collection => {
        return this.renditionConfigurationService
          .forCollection(collection)
          .map(configs => {
            collection.renditionConfigurations = configs;
            return collection;
          });
      })
      .do(c => console.log('collection loader got ', c))
      .subscribe(collection => this.collectionLoaded.emit(collection));
  }

}
