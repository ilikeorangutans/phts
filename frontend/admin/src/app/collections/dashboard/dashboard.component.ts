import { Component, OnInit } from '@angular/core';

import { CollectionService } from "../../services/collection.service";
import { Collection } from "../../models/collection";

@Component({
  selector: 'collections-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  collections: Collection[];

  constructor(
    private collectionService: CollectionService
  ) { }

  ngOnInit() {
    this.collectionService.recent().then((collections) => {
      this.collections = collections;
    });
  }

}
