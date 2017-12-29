import { Component, OnInit } from '@angular/core';

import { CollectionService } from './../../services/collection.service';
import { Collection } from '../../models/collection';

@Component({
  selector: 'collection-form',
  templateUrl: './form.component.html',
  styleUrls: ['./form.component.css']
})
export class FormComponent implements OnInit {

  collection = new Collection();

  constructor(
    private collectionService: CollectionService
  ) { }

  ngOnInit() {
  }

  onSubmit() {
    this.collectionService.save(this.collection);
    this.collection = new Collection();
  }
}
