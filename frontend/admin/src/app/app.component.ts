import { Component, Inject, OnInit } from '@angular/core';
import { CollectionService } from "./services/collection.service";
import { Collection } from "./models/collection";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  constructor(
    private collectionService: CollectionService,
  ) {}

  title = 'app';

  collections: Collection[];

  ngOnInit(): void {
    console.log("AppComponent::ngOnInit()");
    this.collectionService.recent().then((c) => { this.collections = c; });
  }

}
