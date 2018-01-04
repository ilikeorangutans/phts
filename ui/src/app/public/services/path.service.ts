import { DOCUMENT } from '@angular/common';
import { Injectable, Inject } from '@angular/core';

@Injectable()
export class PathService {

  constructor(
    @Inject(DOCUMENT) private document: any,
  ) { }

  apiBase(): string {
    return new URL('/api/', 'http://localhost:8080').toString();
  }

  shareBase(): string {
    return new URL('share/', this.apiBase()).toString();
  }

  shareBySlug(slug: string): string {
    return new URL(slug, this.shareBase()).toString();
  }

  renditionBySlug(slug: string, renditionID: number): string {
    return [this.shareBySlug(slug), 'renditions', renditionID.toString()].join('/');
  }

}
