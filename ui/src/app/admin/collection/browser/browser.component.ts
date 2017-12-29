import { Collection } from './../../models/collection';
import { Component, Input, OnInit } from '@angular/core';

import { CollectionService } from './../../services/collection.service';

@Component({
  selector: 'collection-browser',
  templateUrl: './browser.component.html',
  styleUrls: ['./browser.component.css']
})
export class BrowserComponent implements OnInit {

  @Input()
  numEntries: number = 20;

  collections: Array<Collection> = [];

  constructor(
    private collectionService: CollectionService
  ) { }

  ngOnInit() {
    this.collectionService.recent().then(collections => this.collections = collections);
  }

}
