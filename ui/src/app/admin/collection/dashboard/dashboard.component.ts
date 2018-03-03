import { Collection } from './../../models/collection';
import { Observable } from 'rxjs/Observable';
import { CollectionStore } from './../../stores/collection.store';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  recentCollections: Observable<Array<Collection>>;

  constructor(
    private collectionStore: CollectionStore
  ) { }

  ngOnInit() {
    this.recentCollections = this.collectionStore.recent;
    this.collectionStore.refreshRecent();
  }

}
