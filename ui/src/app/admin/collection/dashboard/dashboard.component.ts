import { Component, OnInit } from '@angular/core';

import { CollectionService } from './../../services/collection.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  constructor(
    private collectionService: CollectionService
  ) { }

  ngOnInit() {
    this.collectionService.setCurrent(null);
  }

}
