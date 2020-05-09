import { Injectable, Inject, isDevMode } from '@angular/core';
import { DOCUMENT } from '@angular/common';

@Injectable({
  providedIn: 'root',
})
export class BasePathService {
  readonly apiHost: string;

  constructor(@Inject(DOCUMENT) private readonly document: Document) {
    if (isDevMode()) {
      this.apiHost = '//localhost:8080';
    } else {
      this.apiHost = `//${new URL(this.document.baseURI).host}`;
    }
    console.log(`base path service initialized with apiHost ${this.apiHost}`);
  }
}
