import { RenditionConfigurationService } from './../../services/rendition-configuration.service';
import { ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { OnDestroy } from '@angular/core';

import { Collection } from './../../models/collection';
import { CollectionService } from './../../services/collection.service';
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-app',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit, OnDestroy {

  currentCollection: Collection = null;

  constructor(
    private collectionService: CollectionService,
    private route: ActivatedRoute
  ) { }

  ngOnInit() {
    this.route.params.subscribe(params => {
      const slug = params['slug'];
      if (slug) {
        this.collectionService.bySlug(slug).then(collection => {
          this.currentCollection = collection;
          this.collectionService.setCurrent(collection);
        });
      } else {
        this.currentCollection = null;
        this.collectionService.setCurrent(null);
      }
    });
  }

  ngOnDestroy(): void {
  }
}
