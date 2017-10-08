import { Collection } from './../../models/collection';
import { CollectionService } from './../../services/collection.service';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-collection-browser',
  templateUrl: './collection-browser.component.html',
  styleUrls: ['./collection-browser.component.css']
})
export class CollectionBrowserComponent implements OnInit {

  collections: Array<Collection> = [];

  constructor(
    private collectionService: CollectionService
  ) { }

  ngOnInit() {
    this.collectionService
      .recent()
      .then((collections) => this.collections = collections);
  }

}
