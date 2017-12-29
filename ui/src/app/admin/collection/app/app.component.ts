import { ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import { OnDestroy } from '@angular/core';

import { Collection } from './../../models/collection';
import { CollectionService } from './../../services/collection.service';
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-app',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit, OnDestroy {

  currentCollection: Collection = null;

  private sub: Subscription;

  constructor(
    private collectionService: CollectionService
  ) { }

  ngOnInit() {
    this.sub = this.collectionService.current.subscribe(collection => {
      console.log('AppComponent got collection change', collection);
      this.currentCollection = collection;
    });

  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }
}
