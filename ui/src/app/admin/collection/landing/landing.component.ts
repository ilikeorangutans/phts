import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import 'rxjs/add/operator/switchMap';

import { CollectionService } from './../../services/collection.service';
import { Collection } from '../../models/collection';

@Component({
  selector: 'app-landing',
  templateUrl: './landing.component.html',
  styleUrls: ['./landing.component.css']
})
export class LandingComponent implements OnInit {

  collection: Collection;

  constructor(
    private route: ActivatedRoute,
    private collectionService: CollectionService
  ) { }

  ngOnInit() {
    this.route.params.switchMap(params => {
      return this.collectionService.bySlug(params['slug']);
    })
    .subscribe(collection => {
      this.collection = collection;
      this.collectionService.setCurrent(collection);
    });
  }

}
