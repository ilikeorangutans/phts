import { Observable } from 'rxjs/Observable';
import { Collection } from './../../models/collection';
import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'collection-browser',
  templateUrl: './browser.component.html',
  styleUrls: ['./browser.component.css']
})
export class BrowserComponent implements OnInit {

  @Input()
  numEntries = 20;

  @Input()
  collections: Observable<Array<Collection>>;

  constructor() { }

  ngOnInit() {
  }
}
