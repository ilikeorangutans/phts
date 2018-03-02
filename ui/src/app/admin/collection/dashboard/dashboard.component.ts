import { CollectionStoreService } from './../../services/collection-store.service';
import { Collection } from './../../models/collection';
import { Observable } from 'rxjs/Observable';
import { Component, OnInit } from '@angular/core';

import { CollectionService } from './../../services/collection.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  recentCollections: Observable<Array<Collection>>;

  constructor(
    private collectionStore: CollectionStoreService
  ) { }

  ngOnInit() {
    this.collectionStore.setCurrent(null);
    this.recentCollections = this.collectionStore.recent();
    this.collectionStore.refreshRecent();
  }

  refresh() {
    this.collectionStore.refreshRecent();
  }
}
