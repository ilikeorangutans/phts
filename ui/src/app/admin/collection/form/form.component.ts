import { CollectionStore } from './../../stores/collection.store';
import { Component, OnInit } from '@angular/core';

import { Collection } from '../../models/collection';

@Component({
  selector: 'collection-form',
  templateUrl: './form.component.html',
  styleUrls: ['./form.component.css']
})
export class FormComponent implements OnInit {

  collection = new Collection();

  constructor(
    private collectionStore: CollectionStore
  ) { }

  ngOnInit() {
  }

  onSubmit() {
    this.collectionStore.save(this.collection);
    this.collection = new Collection();
  }
}
