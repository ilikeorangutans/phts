import { CurrentCollectionService } from './../../collections/current-collection.service';
import { Collection } from './../../models/collection';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-configuration-list',
  templateUrl: './configuration-list.component.html',
  styleUrls: ['./configuration-list.component.css']
})
export class ConfigurationListComponent implements OnInit {

  collection: Collection;

  constructor(
    private currentCollection: CurrentCollectionService
  ) {
    currentCollection.current$.subscribe(collection => {
      this.collection = collection;
      console.log(collection);
    });
  }

  ngOnInit() {
  }

}
