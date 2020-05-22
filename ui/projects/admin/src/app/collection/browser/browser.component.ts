import { Component, Input, OnInit } from '@angular/core';

import { Collection } from './../../models/collection';
import { Observable } from 'rxjs';

@Component({
  selector: 'collection-browser',
  templateUrl: './browser.component.html',
  styleUrls: ['./browser.component.css'],
})
export class BrowserComponent implements OnInit {
  @Input()
  numEntries = 20;

  @Input()
  collections: Observable<Array<Collection>>;

  constructor() {}

  ngOnInit() {}
}
