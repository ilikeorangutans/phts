import { Component, OnInit } from '@angular/core';
import { CollectionService } from '../../services/collection.service';
import { Collection } from '../../models/collection';

@Component({
  selector: 'app-collection-create',
  templateUrl: './collection-create.component.html',
  styleUrls: ['./collection-create.component.css']
})
export class CollectionCreateComponent implements OnInit {

  collection = new Collection();

  constructor(
    private collectionService: CollectionService
  ) { }

  ngOnInit() {
  }

  debug(): String {
    return JSON.stringify(this.collection);
  }

  onSubmit() {
    this.collectionService.save(this.collection);
  }
}
